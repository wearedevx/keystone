#!/bin/bash

if [[ -z "${TMDIR}" ]]; then
    echo "SET TMPDIR"
    export TMPDIR=/tmp/
fi

export $(cat .env-test | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cli/pkg/client.ApiURL=$KSAPI_URL \
	-X github.com/wearedevx/keystone/api/pkg/jwt.salt=${JWT_SALT}"
NOSPIN=true

DBFILE="${TMPDIR}keystone_gorm.db"

if [ -f $DBFILE ]; then
  rm $DBFILE;
fi

# Create db file
touch $DBFILE

cd ../api
make -i run-test &

cd ../cli

echo "go test -tags test -ldflags \"$LDFLAGS\" -work $@"
go test -tags test -ldflags "$LDFLAGS" -work "$@"

EXIT_STATUS_CODE=$?

if [ $EXIT_STATUS_CODE -eq 0 ]; then
	echo "All tests passed";
else
	echo "Some test failed";
fi

if [ -f $DBFILE ]; then
  rm $DBFILE;
fi

# In case the tests failed or succeeded too fast
# the API is not started yet, and lsof fails,
# and the API keeps on running.
# This little for loop here, ensures that
# we wait long enough, ie. when lsof succeeds
for i in {0..10}; do
	pid=$(lsof -t -i :9001);
	if [ $? -eq 0 ]; then
		kill -9 $pid;
		break;
	fi
	sleep 1;
done

exit $EXIT_STATUS_CODE
