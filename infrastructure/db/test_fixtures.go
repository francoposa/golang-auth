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
	// Makes file url driver available to the migrate package
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	"golang-auth/infrastructure/crypto"
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
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

func SetUpAuthNUserRepo(t *testing.T, sqlxDB *sqlx.DB) (repos.AuthNUserRepo, []*resources.AuthNUser) {
	t.Helper()

	authNUserRepo := NewPGAuthNUserRepo(sqlxDB, crypto.NewDefaultArgon2PassHasher())

	users := []*resources.AuthNUser{
		resources.NewAuthNUser("domtoretto", "americanmuscle@fastnfurious.com"),
		resources.NewAuthNUser("brian", "importtuners@fastnfurious.com"),
		resources.NewAuthNUser("roman", "ejectoseat@fastnfurious.com"),
	}

	for _, user := range users {
		_, err := authNUserRepo.Create(user, fmt.Sprintf("%s_pass", user.Username))
		if err != nil {
			panic(err)
		}
	}
	return authNUserRepo, users
}

func SetUpAuthZRoleRepo(t *testing.T, sqlxDB *sqlx.DB) (repos.AuthZRoleRepo, []*resources.AuthZRole) {
	t.Helper()

	repo := NewPGAuthZRoleRepo(sqlxDB)

	roles := []*resources.AuthZRole{
		resources.NewAuthZRole("admin"),
		resources.NewAuthZRole("user"),
	}

	for _, role := range roles {
		_, err := repo.Create(role)
		if err != nil {
			panic(err)
		}
	}
	return repo, roles
}

func SetUpAuthZResourceRepo(t *testing.T, sqlxDB *sqlx.DB) (repos.ResourceRepo, []*resources.AuthZResourceType) {
	t.Helper()

	repo := NewPGAuthZResourceTypeRepo(sqlxDB)

	resourceTypes := []*resources.AuthZResourceType{
		resources.NewAuthZResourceType("user", "ExampleCom User entity"),
		resources.NewAuthZResourceType("profile", "ExampleCom User Profile entity"),
	}

	for _, resourceType := range resourceTypes {
		_, err := repo.Create(resourceType)
		if err != nil {
			panic(err)
		}
	}
	return repo, resourceTypes
}

func SetUpAuthZClientRepo(t *testing.T, sqlxDB *sqlx.DB) (repos.AuthZClientRepo, []*resources.AuthZClient) {
	t.Helper()

	repo := pgAuthZClientRepo{db: sqlxDB}
	clients := []*resources.AuthZClient{}

	for _, client := range clients {
		_, err := repo.Create(client)
		if err != nil {
			panic(err)
		}
	}
	return &repo, clients
}
