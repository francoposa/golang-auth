package domain

import (
	"fmt"
	"strconv"

	validator "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const MinUsernameLen = 8
const MaxUsernameLen = 64
const MinPasswordLen = 16
const MaxPasswordLen = 128

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func NewUser(username, email string) (*User, error) {
	id := uuid.New()

	if !validator.StringLength(
		username,
		strconv.Itoa(MinUsernameLen),
		strconv.Itoa(MaxUsernameLen),
	) {
		return nil, &UsernameInvalidError{}
	}

	if !validator.IsEmail(email) {
		return nil, &EmailInvalidError{Email: email}
	}

	return &User{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

type UsernameInvalidError struct{}

func (e UsernameInvalidError) Error() string {
	return fmt.Sprintf(
		"Username must be between %d and %d characters",
		MinUsernameLen,
		MaxUsernameLen,
	)
}

type EmailInvalidError struct {
	Email string
}

func (e EmailInvalidError) Error() string {
	return fmt.Sprintf("%s is not a valid email address", e.Email)
}

type PasswordInvalidError struct{}

func (e PasswordInvalidError) Error() string {
	return fmt.Sprintf(
		"Password must be between %d and %d characters",
		MinPasswordLen,
		MaxPasswordLen,
	)
}
