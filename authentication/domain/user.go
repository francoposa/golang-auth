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

type Password string

func NewUser(username, email, password string) (*User, *Password, error) {
	id := uuid.NewV4()
	now := time.Now().UTC()

	err := ValidateUsernameRequirements(username)
	if err != nil {
		return nil, nil, err
	}

	validatedPassword, err := ValidatePasswordRequirements(password)
	if err != nil {
		return nil, nil, err
	}

	if !validator.IsEmail(email) {
		return nil, nil, EmailInvalidError{Email: email}
	}

	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}, validatedPassword, nil
}

func ValidateUsernameRequirements(username string) error {
	if !validator.StringLength(
		username,
		strconv.Itoa(MinUsernameLen),
		strconv.Itoa(MaxUsernameLen),
	) {
		return UsernameInvalidError{}
	}
	return nil
}

func ValidatePasswordRequirements(password string) (*Password, error) {
	if !validator.StringLength(
		password,
		strconv.Itoa(MinPasswordLen),
		strconv.Itoa(MaxPasswordLen),
	) {
		return nil, PasswordInvalidError{}
	}
	validatedPassword := Password(password)
	return &validatedPassword, nil
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
