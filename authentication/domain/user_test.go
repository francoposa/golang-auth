package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {

}

func TestValidatePasswordRequirements(t *testing.T) {
	assertions := assert.New(t)

	t.Run("validate password - success", func(t *testing.T) {
		validPassword := string(make([]rune, MinPasswordLen))
		validatedPassword, err := ValidatePasswordRequirements(validPassword)

		assertions.Nil(err)
		expectedPassword := Password(validPassword)
		assertions.Equal(&expectedPassword, validatedPassword)
	})

	t.Run("validate password - invalid length", func(t *testing.T) {
		shortPassword := string(make([]rune, MinPasswordLen-1))
		validatedPassword, err := ValidatePasswordRequirements(shortPassword)

		assertions.Nil(validatedPassword)
		var passLenErr PasswordInvalidLengthError
		assertions.True(errors.As(err, &passLenErr))

		longPassword := string(make([]rune, MaxPasswordLen+1))
		validatedPassword, err = ValidatePasswordRequirements(longPassword)

		assertions.Nil(validatedPassword)
		assertions.True(errors.As(err, &passLenErr))

	})

}
