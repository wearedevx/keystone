#!/bin/sh
export $(cat .env-dev | xargs)

migrate -database=${DATABASE_URL} -path db/migrations $@
