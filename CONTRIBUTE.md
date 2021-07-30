## Development environment

### .env file

```
KSAUTH_URL= # url to the auth cloud function
KSAPI_URL=  # url to the api cloud function
DATABASE_URL= # postgres url for migrations
```

### Development dependencies

1.  Google CLoud CLI: [Recommended Install](https://cloud.google.com/sdk/docs/install)
2.  Cloud SQL Proxy: [Intall it from here](https://cloud.google.com/sql/docs/mysql/sql-proxy#install)
3.  Golang Migrate: `brew install golang-migrate`

### Development Runtime

#### Development database access

```
cloud_sql_proxy \
    -dir /tmp/cloudsql \
    -instances=wearedevx:europe-west6:devx=tcp:5432 \
    -credential_file=functions/ksapi/wearedevx-aa84b56c44de.json
```

**Note :** You need a valid gcloud credentials file

#### Local functions

For the api function

```
cd functions/ksapi/cmd
go run main.go
```

For the auth function

```
cd fucntions/ksapi/cmd
go run main.go
```

### Running the development version

```
./run.sh [options] [command]
```

### Building the release version

```
./build.sh
```

The release binary is `ks`

## Database migration

### Generate migrations

```
./gen.migration.sh create_an_example_table
```

Creates:

- `db/migrations/xxxxxx_create_an_example_table.down.sh`
- `db/migrations/xxxxxx_create_an_example_table.up.sh`

### Run migrations

Upward:

```
./migrate.sh up
```

Backward:

```
./migrate.sh down
```

Specific version:

```
./migrate.sh force $VERSION
```
