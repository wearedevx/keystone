#!/bin/sh
# go clean -modcache

WORK=$PWD
commit=$(git rev-parse HEAD)

function_dir=$WORK/functions

ksauth_dir=$function_dir/ksauth
ksapi_dir=$function_dir/ksapi

echo "Entering ${ksauth_dir}"
cd $ksauth_dir
go get github.com/wearedevx/keystone@$commit
sh ./deploy.sh

echo "Done!\n\n"
cd $WORK

echo "Entering ${ksapi_dir}"
cd $ksapi_dir
go get github.com/wearedevx/keystone@$commit
sh ./deploy.sh

echo "Done!\n\n"
cd $WORK
