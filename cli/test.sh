#!/bin/bash

if [[ -z "${TMDIR}" ]]; then
    echo "SET TMPDIR"
    export TMPDIR=/tmp/
fi


export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cli/pkg/client.ApiURL=$KSAPI_URL"

DBFILE="${TMPDIR}keystone_gorm.db"

# Create db file
touch $DBFILE

cd ../api
go run -tags test main.go &

cd ../cli
# Dpesn't work, with one file name with "-v" param
# # If no test file given, test all files.
# FOLDERTOTEST = $@

# if [ -z "$FOLDERTOTEST" ]; then
#     FOLDERTOTEST="./..."
# fi

# # Start test
# go test -tags test -ldflags "$LDFLAGS" -work "$FOLDERTOTEST"

echo "START TEST"

echo "go test -tags test -ldflags \"$LDFLAGS\" -work $@"
go test -tags test -ldflags "$LDFLAGS" -work "$@"

EXIT_STATUS_CODE=$?

echo "FINISH TEST WITH $EXIT_STATUS_CODE"

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

removeProcessId "keystone_ksapi.pid"

# Delete db file
# echo "rm $DBFILE"
# rm $DBFILE
rm "/tmp/keystone_gorm"*

kill -9 $(lsof -t -i:9001)

exit $EXIT_STATUS_CODE
