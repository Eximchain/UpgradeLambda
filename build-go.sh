#!/bin/bash
rm UpgradeLambda &>/dev/null
GOOS="linux" 
GOARCH="amd64" 
GOPATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/"
. ./build-go-native.sh
mv UpgradeLambda terraform/test105
mv S3BucketWatch terraform/test105
echo "If there are no errors, compiled successfully at `date`!"

