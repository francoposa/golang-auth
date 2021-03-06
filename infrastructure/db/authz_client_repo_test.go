package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-auth/usecases/resources"
)

func TestPGClientRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := SetUpDB(t)
	defer closeDB(t, sqlxDB)
	clientRepo, _ := SetUpAuthZClientRepo(t, sqlxDB)

	publicClient, err := resources.NewClient("telnyx.com", true, true)
	if err != nil {
		t.Error(err)
	}
	_, err = resources.NewClient("spothero.com", false, true)
	if err != nil {
		t.Error(err)
	}

	t.Run("create public client", func(t *testing.T) {
		createdClient, err := clientRepo.Create(publicClient)
		if err != nil {
			t.Error(err)
		}
		assertClient(assertions, publicClient, createdClient)
	})

	t.Run("get public client", func(t *testing.T) {
		retrievedClient, err := clientRepo.Get(publicClient.ID)
		if err != nil {
			t.Error(err)
		}
		assertClient(assertions, publicClient, retrievedClient)
	})
}

func assertClient(a *assert.Assertions, want, got *resources.AuthZClient) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
