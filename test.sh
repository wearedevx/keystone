#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/pkg/client.ksapiURL=$KSAPI_URL"

go test ./tests/... -ldflags "$LDFLAGS" -work  "$@"
