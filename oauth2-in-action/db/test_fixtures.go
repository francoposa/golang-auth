package db

import (
	"testing"

	"github.com/golang-migrate/migrate"
	// Makes postgres driver available to the migrate package
	_ "github.com/golang-migrate/migrate/database/postgres"
	// Makes file url driver available to the migrate package
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/jmoiron/sqlx"
)

var testDBName = "oauth2_in_action_test"

func SetUpDB(t *testing.T) *sqlx.DB {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	migrateUp(t, pgConfig)
	sqlxDb := MustConnect(pgConfig)
	return sqlxDb
}

func SetUpDBData(t *testing.T, sqlxDB *sqlx.DB) {
	t.Helper()
}

func TearDownDB(t *testing.T) {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	migrateDown(t, pgConfig)
}

func migrateUp(t *testing.T, pgConfig PostgresConfig) {
	t.Helper()
	pgURL := BuildConnectionString(pgConfig)
	migration, err := migrate.New("file://oauth2-in-action/db/migrations", pgURL)
	if err != nil {
		panic(err)
	}
	if err := migration.Up(); err != nil {
		panic(err)
	}
}

func migrateDown(t *testing.T, pgConfig PostgresConfig) {
	t.Helper()
	pgURL := BuildConnectionString(pgConfig)
	migration, err := migrate.New("file://oauth2-in-action/db/migrations", pgURL)
	if err != nil {
		panic(err)
	}
	if err := migration.Down(); err != nil {
		panic(err)
	}
}
