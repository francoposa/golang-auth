package db

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func SetUpDB(t *testing.T) *sqlx.DB {
	t.Helper()
	pgConfig := NewDefaultPostgresConfig("oauth2_in_action_test")
	sqlxDb := MustConnect(pgConfig)
	return sqlxDb
}

var truncateTableStatement = "TRUNCATE TABLE %s CASCADE"
var tables = []string{"client"}

func SetUpDBData(t *testing.T, sqlxDB *sqlx.DB) {
	t.Helper()
	truncateTables(t, sqlxDB)
}

func TearDownDBData(t *testing.T, sqlxDB *sqlx.DB) {
	truncateTables(t, sqlxDB)
}

func truncateTables(t *testing.T, sqlxDB *sqlx.DB) {
	t.Helper()
	for _, table := range tables {
		sqlxDB.MustExec(fmt.Sprintf(truncateTableStatement, table))
	}
}
