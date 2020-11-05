package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"golang-auth/authentication/domain"
)

func TestPGAuthNUserRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := SetUpDB(t)
	defer closeDB(t, sqlxDB)
	authNUserRepo, _ := SetUpAuthNUserRepo(t, sqlxDB)

	user, err := domain.NewUser(
		"suki_s2000", "pinkS2000@honda.com",
	)
	assertions.Nil(err)

	t.Run("create authn user", func(t *testing.T) {
		createdAuthNUser, err := authNUserRepo.Create(
			user, "suki_password12345",
		)
		assertions.Nil(err)
		assertions.Equal(user, createdAuthNUser)
	})

	t.Run("create already existing user - error", func(t *testing.T) {
		userWithExistingID := &domain.User{ID: user.ID}
		retrievedUserWithExistingID, err := authNUserRepo.Create(
			userWithExistingID,
			"suki_password12345",
		)
		assertions.Nil(retrievedUserWithExistingID)
		assertions.Equal(
			err,
			domain.UserAlreadyExistsError{
				Field: "id",
				Value: userWithExistingID.ID.String(),
			},
		)

		userWithExistingUsername := &domain.User{Username: user.Username}
		retrievedUserWithExistingUsername, err := authNUserRepo.Create(
			userWithExistingUsername,
			"suki_password12345",
		)
		assertions.Nil(retrievedUserWithExistingUsername)
		assertions.Equal(
			err,
			domain.UserAlreadyExistsError{
				Field: "username",
				Value: userWithExistingUsername.Username,
			},
		)

		userWithExistingEmail := &domain.User{Email: user.Email}
		retrievedUserWithExistingEmail, err := authNUserRepo.Create(
			userWithExistingEmail, "suki_password12345",
		)
		assertions.Nil(retrievedUserWithExistingEmail)
		assertions.Equal(
			err,
			domain.UserAlreadyExistsError{
				Field: "email",
				Value: userWithExistingEmail.Email,
			},
		)
	})

	t.Run("get user by id", func(t *testing.T) {
		retrievedUser, err := authNUserRepo.GetByID(user.ID)
		assertions.Nil(err)
		assertions.Equal(user, retrievedUser)
	})

	t.Run("get nonexistent user by id - error", func(t *testing.T) {
		id := uuid.New()
		retreivedUser, err := authNUserRepo.GetByID(id)
		assertions.Nil(retreivedUser)
		assertions.IsType(domain.UserNotFoundError{}, err)
		assertions.Equal(
			err,
			domain.UserNotFoundError{
				Field: "id",
				Value: id.String(),
			},
		)
	})

	t.Run("get user by username", func(t *testing.T) {
		retrievedUser, err := authNUserRepo.GetByUsername(user.Username)
		assertions.Nil(err)
		assertions.Equal(user, retrievedUser)
	})

	t.Run("get nonexistent user by username - error", func(t *testing.T) {
		retreivedUser, err := authNUserRepo.GetByUsername("xxx")
		assertions.Nil(retreivedUser)
		assertions.IsType(domain.UserNotFoundError{}, err)
		assertions.Equal(
			err,
			domain.UserNotFoundError{
				Field: "username",
				Value: "xxx",
			},
		)
	})

	t.Run("verify authn user password", func(t *testing.T) {
		verified, err := authNUserRepo.VerifyPassword(user.Username, "suki_password12345")
		assertions.Nil(err)
		assertions.True(verified)

		verified, err = authNUserRepo.VerifyPassword(user.Username, "Suki_pass")
		assertions.Nil(err)
		assertions.False(verified)
	})
}
