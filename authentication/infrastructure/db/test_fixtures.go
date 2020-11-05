package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	// Makes postgres driver available to the migrate package
	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/google/uuid"

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

func SetUpDB(t *testing.T) (*sqlx.DB, func(t *testing.T, sqlxDB *sqlx.DB)) {
	t.Helper()

	// Connect to Postgres with user that can create DBs
	pgSuperUserConfig := NewDefaultPostgresConfig("postgres")
	superUserSqlxDB := MustConnect(pgSuperUserConfig)

	// Create random database name to avoid collisions in parallel tests
	rand.Seed(time.Now().UnixNano())
	testDBName := fmt.Sprintf(testDBNameTemplate, rand.Int())
	fmt.Printf("\nCreating test DB %s...\n", testDBName)

	createDBStatement := fmt.Sprintf(createDBStatementTemplate, testDBName)
	superUserSqlxDB.MustExec(createDBStatement)

	// Done with the Postgres superuser - close connection
	err := superUserSqlxDB.Close()
	if err != nil {
		log.Print(err)
	}

	// Connect to test DB with Golang sql db package, as Goose migrations don't work with sqlx.DB
	pgTestDBConfig := NewDefaultPostgresConfig(testDBName)
	pgTestDBURL := BuildConnectionString(pgTestDBConfig)

	fmt.Printf("\nOpening test DB %s...\n", testDBName)
	testDB, err := sql.Open("postgres", pgTestDBURL)
	if err != nil {
		panic(err)
	}

	// Goose migration
	migrateUp(t, testDB)

	// Wrap existing test DB connection into sqlx.DB
	sqlxTestDB := sqlx.NewDb(testDB, "postgres")

	return sqlxTestDB, TearDownDB
}

func TearDownDB(t *testing.T, sqlxDB *sqlx.DB) {
	t.Helper()

	// Extract test DB AuthZResourceType
	var testDBName string
	row := sqlxDB.QueryRowx(`SELECT current_catalog;`)
	err := row.Scan(&testDBName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nClosing test DB %s...\n", testDBName)
	err = sqlxDB.Close()
	if err != nil {
		log.Print(err)
	}

	// Connect to Postgres with user that can drop DBs
	pgSuperUserConfig := NewDefaultPostgresConfig("postgres")
	superUserSqlxDB := MustConnect(pgSuperUserConfig)

	fmt.Printf("\nDropping test DB %s...\n", testDBName)
	createDBStatement := fmt.Sprintf(dropDBStatementTemplate, testDBName)
	superUserSqlxDB.MustExec(createDBStatement)

	// Done with the Postgres superuser - close connection
	err = superUserSqlxDB.Close()
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

func SetUpAuthNUserRepo(t *testing.T, sqlxDB *sqlx.DB) (domain.UserRepo, []*domain.User) {
	t.Helper()

	authNUserRepo := NewPGUserRepo(sqlxDB, crypto.NewDefaultArgon2PassHasher())

	users := []*domain.User{
		{
			ID:       uuid.New(),
			Username: "domtoretto",
			Email:    "americanmuscle@fastnfurious.com",
		},
		{
			ID:       uuid.New(),
			Username: "brian",
			Email:    "importtuners@fastnfurious.com",
		},
		{
			ID:       uuid.New(),
			Username: "roman",
			Email:    "ejectoseat@fastnfurious.com",
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
