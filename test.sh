#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

# Create db file
touch $TMPDIR/keystone_gorm.db

# Start test
go test -tags test -ldflags "$LDFLAGS" -work "$@"

function removeProcessId() {
    kspidfile=$1
    ksapidpath=${TMPDIR}${kspidfile}

    # Check if gcloud auth func pid exist
    if [ -f "$ksapidpath" ]; then
        pid=$(cat $ksapidpath)
        
        echo "kill $kspidfile, PID=$pid"

        rm $ksapidpath

        kill -- -$pid
    fi
}

removeProcessId "keystone_ksauth.pid"
removeProcessId "keystone_ksapi.pid"

# Delete db file
rm $TMPDIR/keystone_gorm.db
