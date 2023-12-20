#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset


# Create Neo4j Helm chart values file
kubectl create ns neo4j
helm repo add neo4j https://helm.neo4j.com/neo4j
helm repo update

# Modify default values like intial DB password
helm install server-1 neo4j/neo4j --namespace neo4j -f neo4j-helm-values.yaml


DEST_DB_KUBERNETES_CONTEXT=koh-aks
DEST_DB_KUBERNETES_NAMESPACE=neo4j
DEST_DB_POD_NAME=server-1-0
DEST_DB_NAME=neo4j
DEST_DB_USER_NAME=neo4j
DEST_DB_PASS=YOUR_PASSWORD # This can be sourced from env file

# Import Environment Nodes from environment.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///environment.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MERGE (e:Environment {name: row.name, description: row.description, id_only_for_import: row.id});"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Workspace Nodes from workspace.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///workspace.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MERGE (w:Workspace {name: row.name, description: row.description, id_only_for_import: row.id});"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Application Nodes from application.csv
CYPHER_QUERY="MATCH (ws:Workspace {name: "kaizentm"}) 
              LOAD CSV WITH HEADERS FROM 'file:///application.csv' AS row
              WITH ws, row WHERE row.id IS NOT NULL
              CREATE (ws)-[:CONTAINS_APP]->(app:Application {name: row.name, description: row.description, id_only_for_import: row.id});"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Workload Nodes from workload.csv
CYPHER_QUERY="MATCH (app:Application)
              LOAD CSV WITH HEADERS FROM 'file:///workload.csv' AS row
              WITH app, row WHERE row.id IS NOT NULL AND row.application_id = app.id_only_for_import
              CREATE (app)-[:CONTAINS_WORKLOAD]->(w:Workload {name: row.name, description: row.description, sourceStorageType: row.source_storage_type, sourceEndpoint: row.source_endpoint, id_only_for_import: row.id});"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}
  

# Import Workload Version Nodes from workload_version.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///workload_version.csv' AS row
              with row where row.id is not null
              MATCH (workload:Workload {id_only_for_import: row.workload_id})
              CREATE (workload)-[:HAS_VERSION]->(:WORKLOAD_VERSION {version: row.version, buildId: row.build_id, sourceCommitId: row.build_commit_id, id_only_for_import: row.id, createdOn: datetime(REPLACE(row.created_on, ' ', 'T'))});"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Deployment Target Nodes from deployment_target.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///deployment_target.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MATCH (e:Environment {id_only_for_import: row.environment_id})
              MATCH (w:Workload {id_only_for_import: row.workload_id})
              CREATE (e)-[:HAS_DEPLOYMENT_TARGET]->(dt:DeploymentTarget {name: row.name, description: row.description, labels: row.labels, manifestStorageType: row.manifests_storage_type, manifestEndpoint: row.manifests_endpoint,  id_only_for_import: row.id})<-[:HAS_DEPLOYMENT_TARGET]-(w);"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Deployment Assignment Nodes from deployment_assignment.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///deployment_assignment.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MATCH (dt:DeploymentTarget {id_only_for_import: row.deployment_target_id})
              MATCH (wv:WORKLOAD_VERSION {id_only_for_import: row.workload_version_id})
              CREATE (dt)-[:HAS_DEPLOYMENT_ASSIGNMENT {gitOpsCommitId: row.gitops_commit_id , createdOn: datetime(REPLACE(row.created_on, ' ', 'T')) ,id_only_for_import: row.id}]->(wv);"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Host Nodes from host.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///host.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MERGE (:Host {name: row.name, description: row.description, labels: (CASE WHEN row.Team IS NULL THEN '' ELSE row .labels END), id_only_for_import: row.id});"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

# Import Reconciler Nodes from reconciler.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///reconciler.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MATCH (h:Host {id_only_for_import: row.host_id})
              MATCH (dt:DeploymentTarget {manifestEndpoint: row.manifests_endpoint})
              CREATE (h)-[:HAS_RECONCILER]->(:Reconciler {name: row.name, description: row.description, reconcilerType: row.reconciler_type, labels: (CASE WHEN row.Team IS NULL THEN '' ELSE row .labels END), manifestStorageType: row.manifests_storage_type, manifestEndpoint: row.manifests_endpoint ,id_only_for_import: row.id})<-[:HAS_RECONCILER]-(dt);"
    
    
kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}


# Import Deployment Nodes from deployment.csv
CYPHER_QUERY="LOAD CSV WITH HEADERS FROM 'file:///deployment.csv' AS row
              WITH row WHERE row.id IS NOT NULL
              MATCH (r: Reconciler {id_only_for_import: row.reconciler_id})
              CREATE (r)-[:HAS_DEPLOYMENT]->(:Deployment {name: row.name, description: row.description, gitOpsCommitId: row.gitops_commit_id, status: row.status, statusMessage: row.status_message , deployedOn: row.created_on, id_only_for_import: row.id})"

kubectl exec --context ${DEST_DB_KUBERNETES_CONTEXT} -n ${DEST_DB_KUBERNETES_NAMESPACE} -it ${DEST_DB_POD_NAME} -- cypher-shell -u ${DEST_DB_USER_NAME} -p ${DEST_DB_PASS} -d ${DEST_DB_NAME} --format plain ${CYPHER_QUERY}

