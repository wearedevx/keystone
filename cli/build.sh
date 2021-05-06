#!/bin/bash

export $(cat .env | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cli/pkg/client.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/cli/pkg/client.ksapiURL=$KSAPI_URL"

go build -ldflags "$LDFLAGS" -o ks
