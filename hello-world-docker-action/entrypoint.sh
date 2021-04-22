#!/bin/sh -l

echo "Hello $1"
echo "ls"
time=$(date)
echo "::set-output name=time::$time"
