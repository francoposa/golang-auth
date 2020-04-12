package db

import (
	"reflect"
	"testing"

	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
)

func TestPGClientRepo(t *testing.T) {
	sqlxDB := SetUpDB(t)
	clientRepo, stubClients := SetUpClientRepo(t, sqlxDB)

	t.Run("create client", func(t *testing.T) {
		clientExample := resources.NewClient("example.com")
		createdClientExample, _ := clientRepo.Create(clientExample)
		assertClient(t, clientExample, createdClientExample)
	})

	t.Run("get created client", func(t *testing.T) {
		stubClient := stubClients[0]
		retrievedClientExample, err := clientRepo.Get(stubClient.ID)
		if err != nil {
			t.Error(err)
		}
		assertClient(t, stubClient, retrievedClientExample)
	})

	TearDownDB(t)
}

func assertClient(t *testing.T, want, got *resources.Client) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nincorrect client resource\nwant: %q, got: %q", want, got)
	}
}
