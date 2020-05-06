package db

import (
	"database/sql"
	"fmt"

	"log"
	"path/filepath"
	"runtime"
	"testing"

	// Makes postgres driver available to the migrate package
	_ "github.com/golang-migrate/migrate/database/postgres"
	// Makes file url driver available to the migrate package
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

var testDBName = "golang_auth_test"

var tables = []string{"auth_user", "client"}

func SetUpDB(t *testing.T) (*sqlx.DB, func(t *testing.T, sqlxDB *sqlx.DB)) {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	pgURL := BuildConnectionString(pgConfig)

	fmt.Print("\nOpening db...\n")
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		panic(err)
	}

	migrateUp(t, db)

	sqlxDB := sqlx.NewDb(db, "postgres")

	fmt.Print("\nTruncating tables for setup...\n\n")
	for _, table := range tables {
		statement := fmt.Sprintf(`TRUNCATE TABLE %s CASCADE`, table)
		result, err := sqlxDB.Queryx(statement)
		if err != nil {
			panic(err)
		}
		query := fmt.Sprintf(`SELECT COUNT(*) FROM %s CASCADE`, table)
		result, err = sqlxDB.Queryx(query)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}

	return sqlxDB, TearDownDB
}

func TearDownDB(t *testing.T, sqlxDB *sqlx.DB) {
	t.Helper()

	fmt.Print("\nTruncating tables for teardown...\n")
	for _, table := range tables {
		statement := fmt.Sprintf(`TRUNCATE TABLE %s CASCADE`, table)
		result, err := sqlxDB.Queryx(statement)
		if err != nil {
			panic(err)
		}
		query := fmt.Sprintf(`SELECT COUNT(*) FROM %s CASCADE`, table)
		result, err = sqlxDB.Queryx(query)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}

	fmt.Print("\nClosing db...\n")
	err := sqlxDB.Close()
	if err != nil {
		log.Print(err)
	}
}

func migrateUp(t *testing.T, db *sql.DB) {
	t.Helper()

	_, dbTestFixturesPath, _, _ := runtime.Caller(1)
	dbPath := filepath.Dir(dbTestFixturesPath)
	migrationsPath := fmt.Sprintf("/%s/migrations", dbPath)

	err := goose.Up(db, migrationsPath)
	if err != nil && err != goose.ErrNoNextVersion {
		panic(err)
	}
}
