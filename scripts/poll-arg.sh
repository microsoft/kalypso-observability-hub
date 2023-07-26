#!/bin/bash

# Polls ARG with the specified interval until the deployment is complete or the timeout is reached.
# A temporary solution before ARG controller
# 
# Usage:
#   poll-arg.sh flags

# Flags:
#   -r       GitOps Repository URL (e.g. https://github.com/microsoft/kalypso-gitops)
#   -b       Environment branch (e.g. dev)
#   -g       Resource Group (e.g. kalypso-dev)

# Example:
#   poll-arg.sh -r https://github.com/microsoft/kalypso-gitops -b dev -g kalypso-dev



while getopts "r:b:g:" option;
    do
    case "$option" in
        r ) REPO_URL=${OPTARG};;
        b ) REPO_BRANCH=${OPTARG};;
        g ) RESOURCE_GROUP=${OPTARG};;
    esac
done

poll_interval=5 # seconds
pg_host="localhost"
pg_user="creator"
pg_database="hub"
reconciler_id=1

export PGPASSWORD=c67

set -eo pipefail  # fail on error
az extension add --name resource-graph

error() {
   echo $1>&2
   exit 1
}

usage() {
echo $1>&2    
cat <<EOM
Usage:
  poll-arg.sh flags

Flags:
  -r       GitOps Repository URL (e.g. https://github.com/microsoft/kalypso-gitops)
  -b       Environment branch (e.g. dev)
  -g       Resource Group (e.g. kalypso-dev)

Example:
  poll-arg.sh -r https://github.com/microsoft/kalypso-gitops -b dev -g kalypso-dev
EOM
exit 1
}

check_parameters() {
    if [ -z $REPO_URL ] && [ -z $REPO_BRANCH ] && [ -z $RESOURCE_GROUP ]
    then
        usage "No arguments specified"
    elif [ -z $REPO_URL ]
    then
        usage "No repository url specified"  
    elif [ -z $REPO_BRANCH ]
    then
        usage "No repository branch specified"  
    elif [ -z $RESOURCE_GROUP ]
    then
        usage "No resource group specified"  
    fi
}

# Queries the Azure Resource Graph for all FluxConfigurations in the specified repository and branch
get_all_configs() {
    total_query="kubernetesconfigurationresources | where type == 'microsoft.kubernetesconfiguration/fluxconfigurations' | where resourceGroup == ""'""$RESOURCE_GROUP""'"" | where properties.gitRepository.url == ""'""$REPO_URL""'"" | where properties.gitRepository.repositoryRef.branch == ""'""$REPO_BRANCH""'"""
    az graph query -q "$total_query"
}


poll() {
attempt=1
while [ true ]
do
    echo "Polling ARG Attempt $attempt ..."
    
    total_configs=$(get_all_configs)

    commit_id=$(echo $total_configs | jq '.data[0].properties.sourceSyncedCommitId')
    compliance_state=$(echo $total_configs | jq '.data[0].properties.complianceState')

    echo $total_configs
    echo $commit_id
    echo $compliance_state
    
    # if $commit_id is not null
    if [ "$commit_id" != "null" ]
    then

      if [ "$compliance_state" == "Compliant" ]
      then
          status="success"
      elif [ "$compliance_state" == "Non-Compliant" ]
      then
          status="failure"
      else
          status="in_progress"
      fi

      query="INSERT INTO deployment (gitops_commit_id, reconciler_id, status, status_message) VALUES ("\'$commit_id\'", "\'$reconciler_id\'", "\'$status\'", '') \
                ON CONFLICT (gitops_commit_id, reconciler_id) DO \
                UPDATE SET status="\'$status\'", \
                          status_message='', \
                          updated_on=current_timestamp, \
                          updated_by=current_user"
      psql -h "$pg_host" -U "$pg_user" -d "$pg_database" -c "$query;"    

    fi




    sleep $poll_interval
    attempt=$(( $attempt + 1 ))

done

}

check_parameters
poll

