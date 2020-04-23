package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	// Makes postgres driver available to the migrate package
	_ "github.com/golang-migrate/migrate/database/postgres"
	// Makes file url driver available to the migrate package
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	"golang-auth/infrastructure/crypto"
	"golang-auth/usecases/interfaces"
	"golang-auth/usecases/resources"
)

var testDBName = "golang_auth_test"

func SetUpDB(t *testing.T) *sqlx.DB {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	migrateDown(t, pgConfig)
	migrateUp(t, pgConfig)
	sqlxDb := MustConnect(pgConfig)
	return sqlxDb
}

func TearDownDB(t *testing.T) {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	migrateDown(t, pgConfig)
}

func migrateUp(t *testing.T, pgConfig PostgresConfig) {
	t.Helper()

	pgURL := BuildConnectionString(pgConfig)

	_, dbTestFixturesPath, _, _ := runtime.Caller(1)
	dbPath := filepath.Dir(dbTestFixturesPath)
	migrationsPath := fmt.Sprintf("/%s/migrations", dbPath)

	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, migrationsPath)
	if err != nil && err == goose.ErrNoNextVersion {
		panic(err)
	}
}

func migrateDown(t *testing.T, pgConfig PostgresConfig) {
	t.Helper()
	pgURL := BuildConnectionString(pgConfig)

	_, dbTestFixturesPath, _, _ := runtime.Caller(1)
	dbPath := filepath.Dir(dbTestFixturesPath)
	migrationsPath := fmt.Sprintf("/%s/migrations", dbPath)

	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		panic(err)
	}

	err = goose.DownTo(db, migrationsPath, 0)
	if err != nil {
		panic(err)
	}
}

func SetUpAuthUserRepo(t *testing.T, sqlxDB *sqlx.DB) (interfaces.AuthUserRepo, []*resources.AuthUser) {
	t.Helper()

	authUserRepo := pgAuthUserRepo{
		db:         sqlxDB,
		passHasher: crypto.NewDefaultArgon2PassHasher(),
	}
	users := []*resources.AuthUser{
		resources.NewAuthUser("domtoretto", "americanmuscle@fastnfurious.com"),
		resources.NewAuthUser("brian", "importtuners@fastnfurious.com"),
		resources.NewAuthUser("roman", "ejectoseat@fastnfurious.com"),
	}

	for _, user := range users {
		_, err := authUserRepo.Create(user, fmt.Sprintf("%s_pass", user.Username))
		if err != nil {
			panic(err)
		}
	}
	return &authUserRepo, users
}

func SetUpClientRepo(t *testing.T, sqlxDB *sqlx.DB) (interfaces.ClientRepo, []*resources.Client) {
	t.Helper()
	clientRepo := pgClientRepo{db: sqlxDB}
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
	return &clientRepo, clients
}
