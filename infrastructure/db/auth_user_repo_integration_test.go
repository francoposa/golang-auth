package db

import (
	"reflect"
	"testing"

	"golang-auth/usecases/resources"
)

func TestPGAuthUserRepo(t *testing.T) {
	sqlxDB := SetUpDB(t)
	authUserRepo, _ := SetUpAuthUserRepo(t, sqlxDB)

	t.Run("create auth user", func(t *testing.T) {
		authUserExample := resources.NewAuthUser("suki", "pink2000@honda.com")
		createdAuthUserExample, _ := authUserRepo.Create(authUserExample, "suki")
		assertAuthUser(t, authUserExample, createdAuthUserExample)
	})

	TearDownDB(t)
}

func assertAuthUser(t *testing.T, want, got *resources.AuthUser) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nincorrect client resource\nwant: %q, got: %q", want, got)
	}
}
