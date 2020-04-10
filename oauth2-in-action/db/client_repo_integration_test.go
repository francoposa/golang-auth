package db

import (
	"reflect"
	"testing"

	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
)

func TestPGClientRepo_Create(t *testing.T) {
	sqlxDB := SetUpDB(t)
	SetUpDBData(t, sqlxDB)

	clientRepo := PGClientRepo{DB: sqlxDB}
	clientExample := resources.NewClient("example.com")
	createdClientExample, _ := clientRepo.Create(clientExample)
	assertClient(t, clientExample, createdClientExample)

	TearDownDBData(t, sqlxDB)
}

func assertClient(t *testing.T, want, got *resources.Client) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nincorrect client resource\nwant: %q, got: %q", want, got)
	}
}
