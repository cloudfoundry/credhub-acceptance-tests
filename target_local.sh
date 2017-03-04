#!/usr/bin/env bash

set -eu

CREDHUB_SRC=${CREDHUB_SRC:-~/workspace/credhub-release/src/credhub}
PKCS12_PATH=${CREDHUB_SRC}/src/test/resources
PKCS12_PASSWORD=${PKCS12_PASSWORD:-changeit}

cat <<EOF > config.json
{
  "api_url": "https://localhost:9000",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "valid_pkcs12_path":"${PKCS12_PATH}/client_cert.p12:${PKCS12_PASSWORD}"
}
EOF
echo "config.json now points to localhost:9000"
