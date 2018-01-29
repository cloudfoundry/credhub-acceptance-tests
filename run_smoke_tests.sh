#!/bin/bash

set -eu

API_URL=${API_URL:-https://localhost:9000}
USERNAME=${USERNAME:-credhub}
PASSWORD=${PASSWORD:-password}
CLIENT_NAME=${CLIENT_NAME:-credhub_client}
CLIENT_SECRET=${CLIENT_SECRET:-secret}

cat <<EOF > test_config.json
{
  "api_url": "${API_URL}",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "client_name":"${CLIENT_NAME}",
  "client_secret":"${CLIENT_SECRET}"
}
EOF

ginkgo -r -p smoke_test
