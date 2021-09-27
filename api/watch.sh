#!/bin/sh
# From https://stackoverflow.com/a/34672970
# Restarts the API on file changes

# Allows for Ctrl-C exit
sigint_handler()
{
  kill $PID
  exit
}

trap sigint_handler SIGINT

while true; do
  # Run the API in a background proccess
  make -i run &
  # Wait for file changes - BLOCKING
  inotifywait -e modify -e move -e create -e delete -e attrib -r `pwd`

  # Some file(s) changed
  # Figure out the pid of the API
  PID=$(lsof -i:9001 -t)
  # try to terminate normally, then force kill if we must:
  kill -TERM $PID || kill -KILL $PID
done
