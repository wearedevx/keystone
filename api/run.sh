#!/bin/sh

SQL_PROXY_PORT=5432
MAX_RETRIES=10
RETRY_DELAY=2

function wait_for_database {
	ready="0"
	tries=0

	echo "Waiting for database to become available...";

	while [ $ready == "0" ]; do
		if [ $tries -eq $MAX_RETRIES ]; then
			echo "Max wait time for database exceeded"
			exit 1;
		fi

		ready="$(lsof -i ":${SQL_PROXY_PORT}" | grep LISTEN | wc -l | xargs)";

		sleep $RETRY_DELAY;

		((tries++));
	done
}

# Start the proxy
./cloud_sql_proxy -instances=$CLOUDSQL_INSTANCE=tcp:$SQL_PROXY_PORT &

# wait for the proxy to spin up
wait_for_database

# Start the server
./server
