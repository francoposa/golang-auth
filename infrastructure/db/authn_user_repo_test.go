package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

func TestPGAuthNUserRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := SetUpDB(t)
	defer closeDB(t, sqlxDB)
	authNUserRepo, _ := SetUpAuthNUserRepo(t, sqlxDB)

	AuthNUser := resources.NewAuthNUser("suki", "pink2000@honda.com")

	t.Run("create authn user", func(t *testing.T) {
		createdAuthNUser, _ := authNUserRepo.Create(AuthNUser, "suki_pass")
		assertAuthNUser(assertions, AuthNUser, createdAuthNUser)
	})

	t.Run("create already existing user - error", func(t *testing.T) {
		alreadyCreatedAuthNUser, err := authNUserRepo.Create(AuthNUser, "suki_pass")
		assertions.Nil(alreadyCreatedAuthNUser, "expected nil struct, got: %q", alreadyCreatedAuthNUser)
		assertions.IsType(repos.AuthNUsernameAlreadyExistsError{}, err)
	})

	t.Run("get authn user", func(t *testing.T) {
		retrievedAuthNUser, err := authNUserRepo.Get(AuthNUser.Username)
		if err != nil {
			t.Error(err)
		}
		assertAuthNUser(assertions, AuthNUser, retrievedAuthNUser)
	})

	t.Run("get nonexistent authn user - error", func(t *testing.T) {
		nonexistentAuthNUser, err := authNUserRepo.Get("xxx")
		assertions.Nil(nonexistentAuthNUser, "expected nil struct, got: %q", nonexistentAuthNUser)
		assertions.IsType(repos.AuthNUserUsernameNotFoundError{}, err)
	})

	t.Run("verify authn user password", func(t *testing.T) {
		verified, err := authNUserRepo.Verify(AuthNUser.Username, "suki_pass")
		if err != nil {
			t.Error(err)
		}
		assertions.True(verified, "correct password was not verified")

		verified, err = authNUserRepo.Verify(AuthNUser.Username, "Suki_pass")
		if err != nil {
			t.Error(err)
		}
		assertions.False(verified, "incorrect password was verified")
	})
}

func assertAuthNUser(a *assert.Assertions, want, got *resources.AuthNUser) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
