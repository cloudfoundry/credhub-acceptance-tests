#!/bin/bash

set -eu

API_URL=${API_URL:-https://localhost:9000}
USERNAME=${USERNAME:-credhub}
PASSWORD=${PASSWORD:-password}
CREDENTIAL_ROOT=${CREDENTIAL_ROOT:-~/workspace/credhub-release/src/credhub/src/test/resources}

cat <<EOF > config.json
{
  "api_url": "${API_URL}",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "credential_root":"${CREDENTIAL_ROOT}"
}
EOF

./generate_certs.py -caKey ${CREDENTIAL_ROOT}/client_ca_private.pem -caCert ${CREDENTIAL_ROOT}/client_ca_cert.pem
ginkgo -r -p
