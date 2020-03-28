package sql

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresConfig defines Postgres SQL connection information
type PostgresConfig struct {
	ApplicationName string        // The name of the application connecting. Useful for attributing sql load.
	Host            string        // The host where the database is located
	Port            uint16        // The port on which the database is listening
	Username        string        // The username for the database
	Password        string        // The password for the database
	Database        string        // The name of the database
	ConnectTimeout  time.Duration // Amount of time to wait before timing out
}

// NewDefaultPostgresConfig creates and return a default postgres configuration.
func NewDefaultPostgresConfig(appName, dbName string) PostgresConfig {
	return PostgresConfig{
		ApplicationName: appName,
		Host:            "localhost",
		Port:            5432,
		Database:        dbName,
		ConnectTimeout:  5 * time.Second,
	}
}

// buildConnectionString transforms the PostgresConfig into a usable connection string for lib/pq.
// If a missing or invalid field is provided, an error is returned.
func buildConnectionString(pc PostgresConfig) (string, error) {
	if pc.Database == "" {
		return "", fmt.Errorf("postgres database name was not specified")
	}
	if pc.ApplicationName == "" {
		return "", fmt.Errorf("application name must be specified to connect to postgres")
	}
	auth := ""
	if pc.Username != "" || pc.Password != "" {
		auth = fmt.Sprintf("%s:%s@", pc.Username, pc.Password)
	}
	url := fmt.Sprintf(
		"postgres://%s%s:%d/%s",
		auth,
		pc.Host,
		pc.Port,
		pc.Database,
	)
	options := []string{fmt.Sprintf("application_name=%s", pc.ApplicationName)}

	if pc.ConnectTimeout.Seconds() > 0 {
		timeoutStr := strconv.Itoa(int(pc.ConnectTimeout.Seconds()))
		options = append(options, fmt.Sprintf("connect_timeout=%s", timeoutStr))
	}

	//TODO figure out SSL
	options = append(options, "sslmode=disable")

	return fmt.Sprintf("%s?%s", url, strings.Join(options, "&")), nil
}

func MustConnect(pc PostgresConfig) *sqlx.DB {
	pgURL, err := buildConnectionString(pc)
	if err != nil {
		panic(err)
	}

	return sqlx.MustConnect("postgres", pgURL)
}
