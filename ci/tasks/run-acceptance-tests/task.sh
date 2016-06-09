#!/bin/bash

set -e -x
echo Hello world
echo $PWD
export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/cm-cli
make dependencies
cd ../cred-hub-acceptance-tests
ginkgo -r integration
