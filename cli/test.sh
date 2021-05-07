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

echo "go test -tags test -ldflags \"$LDFLAGS\" -work $@"
go test -tags test -ldflags "$LDFLAGS" -work "$@"

removeProcessId "keystone_ksapi.pid"

rm "/tmp/keystone_gorm"*

kill -9 $(lsof -t -i:9001)

exit $EXIT_STATUS_CODE
