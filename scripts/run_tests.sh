#!/bin/bash

set -eu

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. && pwd )"

API_URL=${API_URL:-https://localhost:9000}
USERNAME=${USERNAME:-credhub}
PASSWORD=${PASSWORD:-password}
CREDENTIAL_ROOT=${CREDENTIAL_ROOT:-~/workspace/credhub-release/src/credhub/src/test/resources}
UAA_CA=${UAA_CA:-~/workspace/credhub-deployments/ca/uaa_ca.pem}
CLIENT_NAME=${CLIENT_NAME:-credhub_client}
CLIENT_SECRET=${CLIENT_SECRET:-secret}

cat <<EOF > test_config.json
{
  "api_url": "${API_URL}",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}",
  "credential_root":"${CREDENTIAL_ROOT}",
  "uaa_ca":"${UAA_CA}",
  "client_name":"${CLIENT_NAME}",
  "client_secret":"${CLIENT_SECRET}"
}
EOF

pushd "$BASEDIR" >/dev/null
  ginkgo -r -p -skipPackage bbr_integration_test -randomizeAllSpecs -randomizeSuites
popd >/dev/null
