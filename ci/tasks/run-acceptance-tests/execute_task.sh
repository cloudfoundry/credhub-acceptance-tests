#!/bin/bash

set -eu

echo execute_task $PWD

fly \
  -t private \
  execute \
  -c task.yml \
  -i task-repo=../../../ \
  -i subject-repo=../../../../cm-cli
