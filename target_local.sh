#!/usr/bin/env bash

cat <<EOF > config.json
{
  "api_url": "https://localhost:9000",
  "api_username":"${USERNAME}",
  "api_password":"${PASSWORD}"
}
EOF
echo "config.json now points to localhost:9000"
