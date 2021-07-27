#!/bin/sh 

gh secret set DATABASE_URL -b "postgres://keystone:7FMk5pAFdF3aTQt@127.0.0.1:5432/keystone_prod?sslmode=disable"
gh secret set CLOUDSQL_INSTANCE -b "keystone-245200:europe-west6:keystonedb"
gh secret set CLOUDSQL_CREDENTIALS -b "/app/credentials.json"
gh secret set DB_HOST -b "127.0.0.1"
gh secret set DB_NAME -b "keystone_prod"
gh secret set DB_USER -b "keystone"
gh secret set DB_PASSWORD -b "7FMk5pAFdF3aTQt"
gh secret set JWT_SALT -b "qapnx6sP5S9kaeCjZlOmUvgHtsJnA8plRzsQ9whSvnmmBAL922mctk2PYbDzWOj"
gh secret set COMPOSE_PROJECT_NAME -b "keystone"
