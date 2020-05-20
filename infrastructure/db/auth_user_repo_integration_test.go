package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

func TestPGAuthUserRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := SetUpDB(t)
	defer closeDB(t, sqlxDB)
	authUserRepo, _ := SetUpAuthUserRepo(t, sqlxDB)

	authUser := resources.NewAuthUser("suki", "pink2000@honda.com")

	t.Run("create auth user", func(t *testing.T) {
		createdAuthUser, _ := authUserRepo.Create(authUser, "suki_pass")
		assertAuthUser(assertions, authUser, createdAuthUser)
	})

	t.Run("create already existing user - error", func(t *testing.T) {
		alreadyCreatedAuthUser, err := authUserRepo.Create(authUser, "suki_pass")
		assertions.Nil(alreadyCreatedAuthUser, "expected nil struct, got: %q", alreadyCreatedAuthUser)
		assertions.IsType(&repos.DuplicateAuthUserForUsernameError{}, err)
	})

	t.Run("get auth user", func(t *testing.T) {
		retrievedAuthUser, err := authUserRepo.Get(authUser.Username)
		if err != nil {
			t.Error(err)
		}
		assertAuthUser(assertions, authUser, retrievedAuthUser)
	})

	t.Run("get nonexistent auth user - error", func(t *testing.T) {
		nonexistentAuthUser, err := authUserRepo.Get("xxx")
		assertions.Nil(nonexistentAuthUser, "expected nil struct, got: %q", nonexistentAuthUser)
		assertions.IsType(&repos.AuthUserNotFoundForUsernameError{}, err)
	})

	t.Run("verify auth user password", func(t *testing.T) {
		verified, err := authUserRepo.Verify(authUser.Username, "suki_pass")
		if err != nil {
			t.Error(err)
		}
		assertions.True(verified, "correct password was not verified")

		verified, err = authUserRepo.Verify(authUser.Username, "Suki_pass")
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
