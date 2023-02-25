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

	"github.com/francoposa/golang-auth/authentication/domain"
	"github.com/francoposa/golang-auth/authentication/infrastructure/crypto"
)

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
		password := domain.Password(fmt.Sprintf("%s_password12345", user.Username))
		_, err := authNUserRepo.Create(
			user, &password,
		)
		if err != nil {
			panic(err)
		}
	}
	return authNUserRepo, users
}
