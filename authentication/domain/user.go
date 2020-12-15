package domain

import (
	"fmt"
	"strconv"
	"time"

	validator "github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const MinUsernameLen = 8
const MaxUsernameLen = 64
const MinPasswordLen = 16
const MaxPasswordLen = 128

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username, email string) (*User, error) {
	id := uuid.NewV4()

	if !validator.StringLength(
		username,
		strconv.Itoa(MinUsernameLen),
		strconv.Itoa(MaxUsernameLen),
	) {
		return nil, UsernameInvalidError{}
	}

	if !validator.IsEmail(email) {
		return nil, EmailInvalidError{Email: email}
	}

	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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
