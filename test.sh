#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

DBFILE="${TMPDIR}keystone_gorm.db"

# Create db file
touch $DBFILE

# Dpesn't work, with one file name with "-v" param
# # If no test file given, test all files.
# FOLDERTOTEST = $@

# if [ -z "$FOLDERTOTEST" ]; then
#     FOLDERTOTEST="./..."
# fi

# # Start test
# go test -tags test -ldflags "$LDFLAGS" -work "$FOLDERTOTEST"

if [[ -z "${TMDIR}" ]]; then
    export TMPDIR=/tmp
fi

echo "START TEST"

go test -tags test -ldflags "$LDFLAGS" -work "$@"

echo "FINISH TEST"

function removeProcessId() {
    kspidfile=$1
    ksapidpath=${TMPDIR}${kspidfile}

    # Check if gcloud auth func pid exist
    if [ -f "$ksapidpath" ]; then

        echo "rm $kspidfile"

        pid=$(cat $ksapidpath)
        
        echo "kill $kspidfile, PID=$pid"

        rm $ksapidpath

        kill -- -$pid

    else
        echo "File not found: $ksapidpath"
    fi
}

removeProcessId "keystone_ksauth.pid"
removeProcessId "keystone_ksapi.pid"

# Delete db file
echo "rm $DBFILE"
rm $DBFILE
