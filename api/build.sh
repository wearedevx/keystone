#!/bin/sh

docker build -t keystone_server:latest \
	-f Dockerfile .
