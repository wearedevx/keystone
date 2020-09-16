#!/bin/bash

export $(cat .env | xargs)

LDFLAGS="-X github.com/wearedevx/keystone/cmd.ksauthURL=$KSAUTH_URL -X github.com/wearedevx/keystone/cmd.ksapiURL=$KSAPI_URL"

go build -ldflags "$LDFLAGS" -o ks
