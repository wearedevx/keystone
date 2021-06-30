#!/bin/sh
cp "$PWD/snap/snapcraft.template.yaml" "$PWD/snap/snapcraft.yaml"

sed -i "s#<%BRANCH%>#${BRANCH}#g" "$PWD/snap/snapcraft.yaml"
sed -i "s#<%KSAPI_URL%>#${KSAPI_URL}#g" "$PWD/snap/snapcraft.yaml"
sed -i "s#<%GITHUB_CLIENT_ID%>#${GITHUB_CLIENT_ID}#g" "$PWD/snap/snapcraft.yaml"
sed -i "s#<%GITHUB_CLIENT_SECRET%>#${GITHUB_CLIENT_SECRET}#g" "$PWD/snap/snapcraft.yaml"
sed -i "s#<%GITLAB_CLIENT_ID%>#${GITLAB_CLIENT_ID}#g" "$PWD/snap/snapcraft.yaml"
sed -i "s#<%GITLAB_CLIENT_SECRET%>#${GITLAB_CLIENT_SECRET}#g" "$PWD/snap/snapcraft.yaml"
