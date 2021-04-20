#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

#go test -parallel 1 -tags test -ldflags "$LDFLAGS" -work  "$@"

# Create db file
touch $TMPDIR/keystone_gorm.db

# Start test
go test -tags test -ldflags "$LDFLAGS" -work  "$@"

ksauthpidpath=${TMPDIR}keystone_ksauth.pid

pid=`cat $ksauthpidpath`

# Stop gcloud function
# echo "kill -- -$(ps -o pgid=$pid | grep -o '[0-9]*')"
# pkill -P $pid
rm $ksauthpidpath

kill -- -$pid

# Delete db file
rm $TMPDIR/keystone_gorm.db
