package db

import (
	"reflect"
	"testing"

	"golang-auth/usecases/resources"
)

func TestPGAuthUserRepo(t *testing.T) {
	sqlxDB := SetUpDB(t)
	authUserRepo, stubAuthUsers := SetUpAuthUserRepo(t, sqlxDB)

	t.Run("create auth user", func(t *testing.T) {
		authUser := resources.NewAuthUser("suki", "pink2000@honda.com")
		createdAuthUser, _ := authUserRepo.Create(authUser, "suki")
		assertAuthUser(t, authUser, createdAuthUser)
	})

	t.Run("get auth user", func(t *testing.T) {
		stubAuthUser := stubAuthUsers[0]
		retrievedAuthUser, err := authUserRepo.Get(stubAuthUser.ID)
		if err != nil {
			t.Error(err)
		}
		assertAuthUser(t, stubAuthUser, retrievedAuthUser)
	})

	TearDownDB(t)
}

func assertAuthUser(t *testing.T, want, got *resources.AuthUser) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nincorrect client resource\nwant: %q, got: %q", want, got)
	}
}
