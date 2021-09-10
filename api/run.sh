#!/bin/sh

SQL_PROXY_PORT=5432
MAX_RETRIES=10
RETRY_DELAY=2

# Start the proxy
./cloud_sql_proxy \
	-instances=$CLOUDSQL_INSTANCE=tcp:$SQL_PROXY_PORT \
	-credential_file=./credentials.json &

# wait for the proxy to spin up
sleep 10

# Start the server
./server
