#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cli/pkg/client.ApiURL=$KSAPI_URL"

go run -ldflags "$LDFLAGS" main.go "$@"
