#!/bin/bash
rm UpgradeLambda &>/dev/null
echo "GOARCH is: $GOARCH"
echo "GOOS is: $GOOS"
echo "GOPATH is $GOPATH"
go build -o S3BucketWatch -v eximchain.com/S3BucketWatch
go build -o UpgradeLambda -v eximchain.com/UpgradeLambda
go build -o vercmp -v eximchain.com/vercmp

echo "If there are no errors, compiled successfully at `date`!"

