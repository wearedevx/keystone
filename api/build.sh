#!/bin/sh

docker build --no-cache -t keystone_server:latest \
        --build-arg GOOGLE_APPLICATION_CREDENTIALS="keystone-server-credentials.json" \
        --build-arg DB_HOST="$DB_HOST" \
        --build-arg DB_PORT="$DB_PORT" \
        --build-arg DB_NAME="$DB_NAME" \
        --build-arg DB_USER="$DB_USER" \
        --build-arg DB_PASSWORD="$DB_PASSWORD" \
        --build-arg JWT_SALT="$JWT_SALT" \
        --build-arg REDIS_HOST="$REDIS_HOST" \
        --build-arg REDIS_PORT="$REDIS_PORT" \
        --build-arg REDIS_INDEX="$REDIS_INDEX" \
        --build-arg STRIPE_KEY="$STRIPE_KEY" \
        --build-arg STRIPE_WEBHOOK_SECRET="$STRIPE_WEBHOOK_SECRET" \
        --build-arg STRIPE_PRICE="$STRIPE_PRICE" \
        --build-arg X_KS_TTL="$X_KS_TTL" \
	-f Dockerfile .
