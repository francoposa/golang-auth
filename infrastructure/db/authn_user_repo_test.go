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

	AuthNUser := resources.NewAuthNUser("suki", "pinkS2000@honda.com")

	t.Run("create authn user", func(t *testing.T) {
		createdAuthNUser, _ := authNUserRepo.Create(AuthNUser, "suki_pass")
		assertions.Equal(AuthNUser, createdAuthNUser)
	})

	t.Run("create already existing user - error", func(t *testing.T) {
		userWithExistingID := &resources.AuthNUser{ID: AuthNUser.ID}
		retrievedUserWithExistingID, err := authNUserRepo.Create(userWithExistingID, "suki_pass")
		assertions.Nil(retrievedUserWithExistingID)
		assertions.Equal(
			err,
			repos.AuthNUserAlreadyExistsError{
				Field: "id",
				Value: userWithExistingID.ID.String(),
			},
		)

		userWithExistingUsername := &resources.AuthNUser{Username: AuthNUser.Username}
		retrievedUserWithExistingUsername, err := authNUserRepo.Create(userWithExistingUsername, "suki_pass")
		assertions.Nil(retrievedUserWithExistingUsername)
		assertions.Equal(
			err,
			repos.AuthNUserAlreadyExistsError{
				Field: "username",
				Value: userWithExistingUsername.Username,
			},
		)

		userWithExistingEmail := &resources.AuthNUser{Email: AuthNUser.Email}
		retrievedUserWithExistingEmail, err := authNUserRepo.Create(userWithExistingEmail, "suki_pass")
		assertions.Nil(retrievedUserWithExistingEmail)
		assertions.Equal(
			err,
			repos.AuthNUserAlreadyExistsError{
				Field: "email",
				Value: userWithExistingEmail.Email,
			},
		)
	})

	t.Run("get authn user", func(t *testing.T) {
		retrievedAuthNUser, err := authNUserRepo.Get(AuthNUser.Username)
		assertions.Nil(err)
		assertions.Equal(AuthNUser, retrievedAuthNUser)
	})

	t.Run("get nonexistent authn user - error", func(t *testing.T) {
		nonexistentAuthNUser, err := authNUserRepo.Get("xxx")
		assertions.Nil(nonexistentAuthNUser, "expected nil struct, got: %q", nonexistentAuthNUser)
		assertions.IsType(repos.AuthNUserNotFoundError{}, err)
		assertions.Equal(
			err,
			repos.AuthNUserNotFoundError{
				Field: "username",
				Value: "xxx",
			},
		)
	})

	t.Run("verify authn user password", func(t *testing.T) {
		verified, err := authNUserRepo.Verify(AuthNUser.Username, "suki_pass")
		assertions.Nil(err)
		assertions.True(verified, "correct password was not verified")

		verified, err = authNUserRepo.Verify(AuthNUser.Username, "Suki_pass")
		assertions.Nil(err)
		assertions.False(verified, "incorrect password was verified")
	})
}

func assertAuthNUser(a *assert.Assertions, want, got *resources.AuthNUser) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
