#!/bin/bash

export $(cat .env-dev | xargs)
CGO_ENABLED=1

BASE="github.com/wearedevx/keystone/cli"
CLIENT_PKG="${BASE}/pkg/client"
CONSTS_PKG="${BASE}/pkg/constants"
AUTH_PKG="${BASE}/pkg/client/auth"

API_FLAG="-X ${CLIENT_PKG}.ApiURL=$KSAPI_URL"
AUTH_FLAG="-X '${AUTH_PKG}.authRedirectURL=$AUTH_PROXY'"
VERSION_FLAG="-X '${CONSTS_PKG}.Version=$VERSION'"

GITHUB_CLIENT_ID_FLAG="-X ${AUTH_PKG}.githubClientId=$GITHUB_CLIENT_ID"
GITHUB_CLIENT_SECRET_FLAG="-X ${AUTH_PKG}.githubClientSecret=$GITHUB_CLIENT_SECRET"
GITLAB_CLIENT_ID_FLAG="-X ${AUTH_PKG}.gitlabClientId=$GITLAB_CLIENT_ID"
GITLAB_CLIENT_SECRET_FLAG="-X ${AUTH_PKG}.gitlabClientSecret=$GITLAB_CLIENT_SECRET"

LDFLAGS="$API_FLAG \
  $AUTH_FLAG \
  $VERSION_FLAG \
  $GITHUB_CLIENT_ID_FLAG \
  $GITHUB_CLIENT_SECRET_FLAG \
  $GITLAB_CLIENT_ID_FLAG \
  $GITLAB_CLIENT_SECRET_FLAG"

go build -ldflags "$LDFLAGS" -o ks

if [ $? -ne 0 ]; then
  echo "Build failed"; 
else 
  echo "Build done";
fi;
