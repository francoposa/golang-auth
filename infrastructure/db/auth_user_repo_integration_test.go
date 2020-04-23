package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-auth/usecases/resources"
)

func TestPGAuthUserRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB := SetUpDB(t)
	authUserRepo, _ := SetUpAuthUserRepo(t, sqlxDB)

	authUser := resources.NewAuthUser("suki", "pink2000@honda.com")

	t.Run("create auth user", func(t *testing.T) {
		createdAuthUser, _ := authUserRepo.Create(authUser, "suki_pass")
		assertAuthUser(assertions, authUser, createdAuthUser)
	})

	t.Run("get auth user", func(t *testing.T) {
		retrievedAuthUser, err := authUserRepo.Get(authUser.ID)
		if err != nil {
			t.Error(err)
		}
		assertAuthUser(assertions, authUser, retrievedAuthUser)
	})

	t.Run("verify auth user password", func(t *testing.T) {
		verified, err := authUserRepo.Verify(authUser.ID, "suki_pass")
		if err != nil {
			t.Error(err)
		}
		assertions.True(verified, "correct password was not verified")

		verified, err = authUserRepo.Verify(authUser.ID, "Suki_pass")
		if err != nil {
			t.Error(err)
		}
		assertions.False(verified, "incorrect password was verified")
	})
}

func assertAuthUser(a *assert.Assertions, want, got *resources.AuthUser) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
