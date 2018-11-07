#!/bin/bash

set -eu

CLIENT_NAME=${CLIENT_NAME:-credhub_client}
CLIENT_SECRET=${CLIENT_SECRET:-secret}

cat <<EOF > test_config.json
{
  "director_host":"${API_IP}",
  "api_url": "https://${API_IP}:8844",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "bosh": {
    "bosh_environment":"$BOSH_ENVIRONMENT",
    "bosh_client":"$BOSH_CLIENT",
    "bosh_client_secret":"$BOSH_CLIENT_SECRET",
    "bosh_ca_cert_path":"$BOSH_CA_CERT_PATH"
  },
  "uaa_ca":"${SERVER_CA_CERT_PATH}",
  "client_name":"${CLIENT_NAME}",
  "client_secret":"${CLIENT_SECRET}",
  "deployment_name":"$DEPLOYMENT_NAME"
}
EOF

ginkgo -r -p bbr_integration_test
