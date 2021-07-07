#!/bin/sh

migrate create -ext sql -dir db/migrations -seq $@
