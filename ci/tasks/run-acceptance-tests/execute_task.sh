#!/bin/bash

set -eu

echo execute_task $PWD

fly \
  -t private \
  execute \
  --tag vsphere \
  -c task.yml \
  -i task-repo=../../../ \
  -i subject-repo=../../../../credhub-cli
