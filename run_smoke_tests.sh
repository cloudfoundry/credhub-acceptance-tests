#!/bin/bash

set -eu

API_URL=${API_URL:-https://localhost:9000}
USERNAME=${USERNAME:-credhub}
PASSWORD=${PASSWORD:-password}

cat <<EOF > config.json
{
  "api_url": "${API_URL}",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}"
}
EOF

ginkgo -r -p smoke_tests
