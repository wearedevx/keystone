#!/bin/sh

# Start the proxy
./cloud_sql_proxy -instances=$CLOUDSQL_INSTANCE=tcp:5432 &

# wait for the proxy to spin up
sleep 10

# Start the server
./server
