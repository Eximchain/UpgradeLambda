#!/bin/bash
GOOS="linux" 
GOARCH="amd64" 
GOPATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/" 
brew install dep
brew upgrade dep
echo "GOPATH is at: $GOPATH"
pushd src/eximchain.com/UpgradeLambda
dep ensure
popd
pushd src/softwareupgrade
dep ensure
popd
pushd src/eximchain.com/S3BucketWatch
dep ensure
popd

