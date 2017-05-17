#!/bin/bash

set -eu

API_URL=${API_URL:-https://10.244.0.2:8844}
USERNAME=${USERNAME:-credhub_cli}
PASSWORD=${PASSWORD:-credhub_cli_password}
BOSH_URL=${BOSH_URL:-https://192.168.50.4:25555}
BOSH_CLIENT=${BOSH_CLIENT:-admin}
BOSH_CLIENT_SECRET=${BOSH_CLIENT_SECRET:-admin}
BOSH_CERT_PATH=${BOSH_CERT_PATH:-$HOME/workspace/bosh-lite/ca/certs/ca.crt}
CREDHUB_DEPLOYMENT_NAME=${CREDHUB_DEPLOYMENT_NAME:-credhub-lite}

cat <<EOF > config.json
{
  "api_url": "${API_URL}",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "bosh": {
    "url":"${BOSH_URL}",
    "client":"${BOSH_CLIENT}",
    "client_secret":"${BOSH_CLIENT_SECRET}",
    "cert_path":"${BOSH_CERT_PATH}",
    "deployment_name":"${CREDHUB_DEPLOYMENT_NAME}"
  }
}
EOF

ginkgo -r -p bbr_integration_test
