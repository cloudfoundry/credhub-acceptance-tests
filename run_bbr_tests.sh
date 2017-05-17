#!/bin/bash

set -eu

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
