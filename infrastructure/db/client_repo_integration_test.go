package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-auth/usecases/resources"
)

func TestPGClientRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB := SetUpDB(t)
	clientRepo, _ := SetUpClientRepo(t, sqlxDB)

	client := resources.NewClient("example.com")

	t.Run("create client", func(t *testing.T) {
		createdClient, err := clientRepo.Create(client)
		if err != nil {
			t.Error(err)
		}
		assertClient(assertions, client, createdClient)
	})

	t.Run("get client", func(t *testing.T) {
		retrievedClient, err := clientRepo.Get(client.ID)
		if err != nil {
			t.Error(err)
		}
		assertClient(assertions, client, retrievedClient)
	})
}

func assertClient(a *assert.Assertions, want, got *resources.Client) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
