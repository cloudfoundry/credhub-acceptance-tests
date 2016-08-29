#!/bin/bash

set -eu

echo execute_task $PWD

export API_URL=https://50.17.59.67:8844

fly \
  -t private \
  execute \
  --tag vsphere \
  -c task.yml \
  -i task-repo=../../../ \
  -i subject-repo=../../../../credhub-cli
