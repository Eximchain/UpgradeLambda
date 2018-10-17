#!/bin/bash
rm -rf pkg src/eximchain.com/S3BucketWatch/vendor src/eximchain.com/UpgradeLambda/vendor src/github.com src/golang.* src/softwareupgrade/vendor src/gopkg.in src/.DS_Store src/eximchain.com/.DS_Store src/eximchain.com/UpgradeLambda/.DS_Store vercmp S3BucketWatch UpgradeLambda
go clean
