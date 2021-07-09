#!/bin/sh
TARGET="$PWD/snap/snapcraft.yaml"
TEMPLATE="$PWD/cli/deploy/snapcraft.template.yaml"

cp $TEMPLATE $TARGET

sed -i "s#<%BRANCH%>#${BRANCH}#g" $TARGET
sed -i "s#<%VERSION%>#${VERSION}#g" $TARGET

sed -i "s#<%KSAPI_URL%>#${KSAPI_URL}#g" $TARGET
sed -i "s#<%AUTH_PROXY%>#${AUTH_PROXY}#g" $TARGET

sed -i "s#<%GITHUB_CLIENT_ID%>#${GITHUB_CLIENT_ID}#g" $TARGET
sed -i "s#<%GITHUB_CLIENT_SECRET%>#${GITHUB_CLIENT_SECRET}#g" $TARGET
sed -i "s#<%GITLAB_CLIENT_ID%>#${GITLAB_CLIENT_ID}#g" $TARGET
sed -i "s#<%GITLAB_CLIENT_SECRET%>#${GITLAB_CLIENT_SECRET}#g" $TARGET
