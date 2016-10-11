#!/bin/bash

if [ -e config.json ]; then
  echo "config.json already exists, not generating."
  exit 1
else
  cat <<EOF > config.json
{
  "api_url": "http://localhost:9000"
}
EOF
  echo "config.json generated, pointing to localhost:9000"
  exit 0
fi
