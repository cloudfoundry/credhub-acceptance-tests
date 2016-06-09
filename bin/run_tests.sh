#!/bin/bash

set -e -x
echo $PWD
export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

go install -v github.com/onsi/ginkgo/ginkgo
cd $GOPATH/src/github.com/pivotal-cf/cred-hub-acceptance-tests
ginkgo -r integration
