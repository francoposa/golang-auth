package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresConfig defines Postgres SQL connection information
type PostgresConfig struct {
	Host                  string // The host where the database is located
	Port                  uint16 // The port on which the database is listening
	Username              string // The username for the database
	Password              string // The password for the database
	Database              string // The name of the database
	ConnectTimeoutSeconds int    // Number of seconds to wait before timing out
}

// NewDefaultPostgresConfig creates and return a default postgres configuration.
func NewDefaultPostgresConfig(appName, dbName string) PostgresConfig {
	return PostgresConfig{
		Host:                  "localhost",
		Port:                  5432,
		Username:              "postgres",
		Database:              dbName,
		ConnectTimeoutSeconds: 5,
	}
}

// buildConnectionString transforms the PostgresConfig into a usable connection string for lib/pq.
// If a missing or invalid field is provided, an error is returned.
func buildConnectionString(pc PostgresConfig) string {

	auth := ""
	if pc.Username != "" || pc.Password != "" {
		auth = fmt.Sprintf("%s:%s@", pc.Username, pc.Password)
	}
	url := fmt.Sprintf(
		"postgres://%s%s:%d/%s?connect_timeout=%d&sslmode=disable",
		auth,
		pc.Host,
		pc.Port,
		pc.Database,
		pc.ConnectTimeoutSeconds,
	)
	return url
}

func MustConnect(pc PostgresConfig) *sqlx.DB {
	pgURL := buildConnectionString(pc)
	return sqlx.MustConnect("postgres", pgURL)
}
