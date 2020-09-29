#!/bin/bash

export $(cat .env-dev | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cmd.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/cmd.ksapiURL=$KSAPI_URL"

go run -ldflags "$LDFLAGS" main.go "$@"