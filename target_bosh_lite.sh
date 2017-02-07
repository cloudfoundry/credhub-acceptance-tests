#!/bin/bash

cat <<EOF > config.json
{
  "api_url": "https://10.244.0.2:8844",
  "api_username":"credhub_cli",
  "api_password":"credhub_cli_password"
}
EOF
echo "config.json now points to 10.244.0.2:8844"
exit 0
