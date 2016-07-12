#!/bin/bash

set -ex
export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/cm-cli
make dependencies
cd ../cred-hub-acceptance-tests
ginkgo -r integration
