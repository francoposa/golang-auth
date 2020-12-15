package db

import (
	"regexp"

	pgTools "github.com/francoposa/go-tools/postgres"

	"github.com/lib/pq"
)

// NewDefaultPostgresConfig creates and return a default postgres configuration.
func NewDefaultPostgresConfig(dbName string) pgTools.Config {
	return pgTools.Config{
		Host:                  "localhost",
		Port:                  5432,
		Username:              "postgres",
		Password:              "",
		Database:              dbName,
		ApplicationName:       "",
		ConnectTimeoutSeconds: 0,
		SSLMode:               "disable",
	}
}

var pqErrorDetailRegex = regexp.MustCompile(`Key\s\((?P<Key>[a-zA-Z0-9_]*)\)=\((?P<Value>.*)\).*`)

func GetAlreadyExistsErrorKeyValue(err *pq.Error) (string, string) {
	matches := pqErrorDetailRegex.FindStringSubmatch(err.Detail)
	return matches[1], matches[2]
}
