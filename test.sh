#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

# Create db file
touch $TMPDIR/keystone_gorm.db

# Start test
go test -tags test -ldflags "$LDFLAGS" -work "$@"

ksauthpidpath=${TMPDIR}keystone_ksauth.pid

# Check if gcloud auth func pid exist
if [ -f "$ksauthpidpath" ]; then

    pid=$(cat $ksauthpidpath)

    rm $ksauthpidpath

    kill -- -$pid
fi

# Delete db file
rm $TMPDIR/keystone_gorm.db
