#!/bin/bash

URL=$1
DASHBOARDS_DIR=$2

for db_uid in $(curl -s "${URL}/api/search" | jq -r .[].uid); do
  db_json=$(curl -s "${URL}/api/dashboards/uid/${db_uid}")
  db_slug=$(echo "${db_json}" | jq -r .meta.slug)
  db_title=$(echo "${db_json}" | jq -r .dashboard.title)
  filename="${DASHBOARDS_DIR}/${db_slug}.json"
  echo "Exporting \"${db_title}\" to \"${filename}\"..."
  echo "${db_json}" | jq -r .dashboard > "${filename}"
done
