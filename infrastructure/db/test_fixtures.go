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

	"golang-auth/infrastructure/crypto"
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

var testDBName = "golang_auth_test"

func SetUpDB(t *testing.T) (*sqlx.DB, func(t *testing.T, sqlxDB *sqlx.DB)) {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig(testDBName)
	pgURL := BuildConnectionString(pgConfig)

	fmt.Print("\nOpening DB...\n\n")
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		panic(err)
	}

	migrateUp(t, db)

	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, CloseDB
}

func CloseDB(t *testing.T, sqlxDB *sqlx.DB) {
	t.Helper()
	fmt.Print("\nClosing DB...\n\n")
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

func SetUpAuthUserRepo(t *testing.T, sqlxDB *sqlx.DB) (repos.AuthUserRepo, []*resources.AuthUser) {
	t.Helper()

	sqlxDB.MustExec(`TRUNCATE TABLE auth_user CASCADE;`)

	authUserRepo := pgAuthUserRepo{
		db:     sqlxDB,
		hasher: crypto.NewDefaultArgon2PassHasher(),
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

func SetUpClientRepo(t *testing.T, sqlxDB *sqlx.DB) (repos.ClientRepo, []*resources.Client) {
	t.Helper()

	sqlxDB.MustExec(`TRUNCATE TABLE client CASCADE;`)

	clientRepo := pgClientRepo{db: sqlxDB}
	clients := []*resources.Client{}

	for _, client := range clients {
		_, err := clientRepo.Create(client)
		if err != nil {
			panic(err)
		}
	}
	return &clientRepo, clients
}
