#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

go run -ldflags "$LDFLAGS" main.go "$@"
