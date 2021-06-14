#!/bin/sh
# go clean -modcache

WORK=$PWD
commit=$(git rev-parse HEAD)

export $(cat .env | sed 's/#.*//g' | xargs)

TAG=eu.gcr.io/keystone-245200/keystone-server:${commit}

gcloud auth activate-service-account --project=keystone-245200 --key-file=keystone-deploy-credentials.json

gcloud builds submit --tag $TAG

gcloud run deploy keystone-server \
	--region europe-west6 \
	--allow-unauthenticated \
	--set-env-vars DB_HOST=${DB_HOST},DB_NAME=${DB_NAME},DB_USER=${DB_USER},DB_PASSWORD=${DB_PASSWORD},CLOUDSQL_INSTANCE=${CLOUDSQL_INSTANCE},CLOUDSQL_CREDENTIALS=${CLOUDSQL_CREDENTIALS} \
	--add-cloudsql-instances keystonedb \
	--image $TAG

