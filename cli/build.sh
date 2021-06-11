#!/bin/bash

export $(cat .env | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cli/pkg/client.ApiURL=$KSAPI_URL"

go build -ldflags "$LDFLAGS" -o ks
