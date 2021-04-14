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
psql -U postgres -c "CREATE DATABASE golang_auth_authentication;"
```


Run migrations from project root:

```
goose -dir infrastructure/db/migrations postgres "user=postgres dbname=golang_auth_authentication sslmode=disable" up
```

The test setup currently handles running migrations up & down migrations to start & end with an empty database.
