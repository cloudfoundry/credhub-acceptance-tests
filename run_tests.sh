#!/bin/bash

if [ ! -e config.json ]; then
  echo "You are missing config.json."
  exit 1
else
ginkgo -r -p integration_test
fi
