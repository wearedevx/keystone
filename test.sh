#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

DBFILE="${TMPDIR}keystone_gorm.db"

# Create db file
touch $DBFILE

# If no test file given, test all files.
FOLDERTOTEST = $@

if [ -z "$FOLDERTOTEST" ]; then
    FOLDERTOTEST="./..."
fi

# Start test
go test -tags test -ldflags "$LDFLAGS" -work "$FOLDERTOTEST"

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
rm $DBFILE
