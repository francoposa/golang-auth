package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-auth/usecases/resources"
)

func TestPGAuthUserRepo(t *testing.T) {
	assert := assert.New(t)
	sqlxDB := SetUpDB(t)
	authUserRepo, _ := SetUpAuthUserRepo(t, sqlxDB)

	authUser := resources.NewAuthUser("suki", "pink2000@honda.com")

	t.Run("create auth user", func(t *testing.T) {
		createdAuthUser, _ := authUserRepo.Create(authUser, "suki_pass")
		assert.Equal(
			authUser,
			createdAuthUser,
			"expected equivalent structs, want: %q, got: %q",
			authUser,
			createdAuthUser,
		)
	})

	t.Run("get auth user", func(t *testing.T) {
		retrievedAuthUser, err := authUserRepo.Get(authUser.ID)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(
			authUser,
			retrievedAuthUser,
			"expected equivalent structs, want: %q, got: %q",
			authUser,
			retrievedAuthUser,
		)
	})

	t.Run("verify auth user password", func(t *testing.T) {
		verified, err := authUserRepo.Verify(authUser.ID, "suki_pass")
		if err != nil {
			t.Error(err)
		}
		assert.True(verified, "correct password was not verified")

		verified, err = authUserRepo.Verify(authUser.ID, "Suki_pass")
		if err != nil {
			t.Error(err)
		}
		assert.False(verified, "incorrect password was verified")
	})
}
