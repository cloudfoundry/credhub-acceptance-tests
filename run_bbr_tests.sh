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
    "host":"${API_IP}:22",
    "bosh_ssh_username":"${BOSH_SSH_USERNAME}",
    "bosh_ssh_private_key_path":"${BOSH_SSH_PRIVATE_KEY_PATH}"
  },
  "uaa_ca":"${SERVER_CA_CERT_PATH}"
  "client_name":"${CLIENT_NAME}",
  "client_secret":"${CLIENT_SECRET}"
}
EOF

ginkgo -r -p bbr_integration_test
