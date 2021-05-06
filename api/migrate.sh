#!/bin/sh
export $(cat .env | xargs)

migrate -database=${DATABASE_URL} -path db/migrations $@
