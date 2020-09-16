#!/bin/sh
# go clean -modcache

WORK=$PWD
commit=$(git rev-parse HEAD)

cd $PWD/functions/ksauth
go get -u github.com/wearedevx/keystone@$commit
sh ./deploy.sh

cd $WORK

cd $PWD/functions/ksapi
go get -u github.com/wearedevx/keystone@$commit
sh ./deploy.sh

cd $WORK
