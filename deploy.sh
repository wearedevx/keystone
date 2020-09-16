#!/bin/sh

WORK=$PWD

cd $PWD/functions/ksauth
go clean -modcache
go get -u github.com/wearedevx/keystone@go
sh ./deploy.sh

cd $WORK

cd $PWD/functions/ksapi
go clean -modcache
go get -u github.com/wearedevx/keystone@go
sh ./deploy.sh

cd $WORK
