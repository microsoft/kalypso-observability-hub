#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

SOURCE_DB_KUBERNETES_CONTEXT=observability-hub
SOURCE_DB_KUBERNETES_NAMESPACE=hub
SOURCE_DB_POD_NAME=postgre-db-0
SOURCE_DB_NAME=hub
SOURCE_DB_USER_NAME=creator

# Export KOH Postgres tables to CSV files in current directory
# Note: The last line of each CSV file is removed because it contains the number of rows in the table
for TABLE in reconciler deployment deployment_assignment workload_version deployment_target host workspace workload application environment; do
  kubectl exec --context ${SOURCE_DB_KUBERNETES_CONTEXT}  -n ${SOURCE_DB_KUBERNETES_NAMESPACE} ${SOURCE_DB_POD_NAME} -- psql -h localhost -U ${SOURCE_DB_USER_NAME} -d ${SOURCE_DB_NAME} -c "SELECT * FROM ${TABLE};" -A -F ',' | sed '$d' > ${TABLE}.csv
done