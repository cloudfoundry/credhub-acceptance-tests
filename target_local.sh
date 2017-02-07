#!/bin/bash

cat <<EOF > config.json
{
  "api_url": "http://localhost:9000",
  "api_username":"credhub_cli",
  "api_password":"silent42opportunity"
}
EOF
echo "config.json now points to localhost:9000"
exit 0
