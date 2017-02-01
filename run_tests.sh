#!/bin/bash

if [ ! -e config.json ]; then
  echo "You are missing config.json - use ./target_local.sh to point to your local Credhub."
  exit 1
else
ginkgo -r -p integration_test
fi
