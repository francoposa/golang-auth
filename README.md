# golang-auth

Authentication & Authorization with Golang & Postgres

### DB Migrations

Uses [Goose](https://github.com/pressly/goose) for migrations.

Install Goose:

```
$ go get -u github.com/pressly/goose/cmd/goose
```

Create DBs
```
$ psql -U postgres -c "CREATE DATABASE examplecom_auth;"
$ psql -U postgres -c "CREATE DATABASE examplecom_auth_test;"
```


Run migrations from project root:

```
$ goose -dir infrastructure/db/migrations postgres "user=postgres dbname=examplecom_auth sslmode=disable" up
```

The test setup currently handles running migrations up & down migrations to start & end with an empty database.
