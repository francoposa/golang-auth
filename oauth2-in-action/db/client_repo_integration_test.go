package db

import (
	"reflect"
	"testing"

	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
)

func TestPGClientRepo(t *testing.T) {
	sqlxDB := SetUpDB(t)
	clientRepo := PGClientRepo{DB: sqlxDB}

	clientExample := resources.NewClient("example.com")

	t.Run("create client", func(t *testing.T) {
		createdClientExample, _ := clientRepo.Create(clientExample)
		assertClient(t, clientExample, createdClientExample)
	})

	t.Run("get created client", func(t *testing.T) {
		retrievedClientExample, err := clientRepo.Get(clientExample.ID)
		if err != nil {
			t.Error(err)
		}
		assertClient(t, clientExample, retrievedClientExample)
	})

	TearDownDB(t)
}

func assertClient(t *testing.T, want, got *resources.Client) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nincorrect client resource\nwant: %q, got: %q", want, got)
	}
}
