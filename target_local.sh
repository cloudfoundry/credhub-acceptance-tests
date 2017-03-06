#!/usr/bin/env bash

set -eu

CREDHUB_SRC=${CREDHUB_SRC:-~/workspace/credhub-release/src/credhub}
CERTS_PATH=${CREDHUB_SRC}/src/test/resources

cat <<EOF > config.json
{
  "api_url": "https://localhost:9000",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "valid_cert_path":"${CERTS_PATH}/client.pem",
  "valid_private_key_path":"${CERTS_PATH}/client_key.pem"
}
EOF
echo "config.json now points to localhost:9000"
