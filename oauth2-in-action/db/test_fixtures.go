package db

import (
	"testing"

	"github.com/golang-migrate/migrate"
	// Makes postgres driver available to the migrate package
	_ "github.com/golang-migrate/migrate/database/postgres"
	// Makes file url driver available to the migrate package
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/jmoiron/sqlx"

	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
)

var testDBName = "oauth2_in_action_test"

func SetUpDB(t *testing.T) *sqlx.DB {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	migrateUp(t, pgConfig)
	sqlxDb := MustConnect(pgConfig)
	return sqlxDb
}

func TearDownDB(t *testing.T) {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	migrateDown(t, pgConfig)
}

var noChangeErr = "no change"

func migrateUp(t *testing.T, pgConfig PostgresConfig) {
	t.Helper()
	pgURL := BuildConnectionString(pgConfig)
	migration, err := migrate.New("file://db/migrations", pgURL)
	if err != nil && err.Error() != noChangeErr {
		panic(err)
	}
	err = migration.Up()
	if err != nil && err.Error() != noChangeErr {
		panic(err)
	}
}

func migrateDown(t *testing.T, pgConfig PostgresConfig) {
	t.Helper()
	pgURL := BuildConnectionString(pgConfig)
	migration, err := migrate.New("file://db/migrations", pgURL)
	if err != nil {
		panic(err)
	}
	if err := migration.Down(); err != nil {
		panic(err)
	}
}

func SetUpClientRepo(t *testing.T, sqlxDB *sqlx.DB) (PGClientRepo, []*resources.Client) {
	t.Helper()
	clientRepo := PGClientRepo{DB: sqlxDB}
	clients := []*resources.Client{
		resources.NewClient("qualtrics.com"),
		resources.NewClient("telnyx.com"),
		resources.NewClient("spothero.com"),
	}

	for _, client := range clients {
		_, err := clientRepo.Create(client)
		if err != nil {
			panic(err)
		}
	}
	return clientRepo, clients
}
