package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	pgTools "github.com/francoposa/go-tools/postgres"
	sqlTools "github.com/francoposa/go-tools/postgres/database_sql"

	// Makes postgres driver available to the migrate package
	_ "github.com/golang-migrate/migrate/database/postgres"
	uuid "github.com/satori/go.uuid"

	// Makes file url driver available to the migrate package
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	"golang-auth/authentication/domain"
	"golang-auth/authentication/infrastructure/crypto"
)

var testDBNameTemplate = `examplecom_auth_test_%d`
var createDBStatementTemplate = `CREATE DATABASE %s`
var dropDBStatementTemplate = `DROP DATABASE %s`

func SetUpDB(t *testing.T, dbName string, superUserPGConfig pgTools.Config) (*sql.DB, error) {
	t.Helper()

	testDB, err := sqlTools.CreateDB(t, dbName, superUserPGConfig)
	if err != nil {
		return testDB, err
	}

	err = MigrateUp(t, testDB)
	if err != nil {
		return testDB, err
	}

	return testDB, nil
}

func MigrateUp(t *testing.T, db *sql.DB) error {
	t.Helper()

	// Golang is a PITA to pin down working directory during tests so use this
	// https://stackoverflow.com/questions/23847003/golang-tests-and-working-directory
	_, dbTestFixturesPath, _, _ := runtime.Caller(1)
	dbPath := filepath.Dir(dbTestFixturesPath)
	migrationsPath := fmt.Sprintf("/%s/migrations", dbPath)

	err := goose.Up(db, migrationsPath)
	if err != nil && err != goose.ErrNoNextVersion {
		return err
	}

	return nil
}

func SetUpUserRepo(t *testing.T, sqlxDB *sqlx.DB) (domain.UserRepo, []*domain.User) {
	t.Helper()

	authNUserRepo := NewPGUserRepo(sqlxDB, crypto.NewDefaultArgon2PassHasher())

	users := []*domain.User{
		{
			ID:        uuid.NewV4(),
			Username:  "domtoretto",
			Email:     "americanmuscle@fastnfurious.com",
			Enabled:   true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        uuid.NewV4(),
			Username:  "brian",
			Email:     "importtuners@fastnfurious.com",
			Enabled:   true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        uuid.NewV4(),
			Username:  "roman",
			Email:     "ejectoseat@fastnfurious.com",
			Enabled:   true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	for _, user := range users {
		_, err := authNUserRepo.Create(
			user, fmt.Sprintf("%s_password12345", user.Username),
		)
		if err != nil {
			panic(err)
		}
	}
	return authNUserRepo, users
}
