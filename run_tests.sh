#!/bin/bash

set -eu

USERNAME=${USERNAME:-credhub}
PASSWORD=${PASSWORD:-password}
CREDENTIAL_ROOT=${CREDHUB_SRC:-~/workspace/credhub-release/src/credhub/src/test/resources}
API_URL=${API_URL:-https://localhost:9000}

cat <<EOF > config.json
{
  "api_url": "${API_URL}",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "credential_root":"${CREDENTIAL_ROOT}"
}
EOF

./generate_certs.py -caKey $CREDENTIAL_ROOT/ca_key.pem -caCert $CREDENTIAL_ROOT/ca.pem
ginkgo -r -p
