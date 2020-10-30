package domain

import (
	"fmt"
	"regexp"
)

type EmailAddress struct {
	Email string
}

func NewEmailAddress(email string) (*EmailAddress, error) {
	if !isEmailValid(email) {
		return nil, &EmailInvalidError{Email: email}
	}
	return &EmailAddress{Email: email}, nil
}

func (e EmailAddress) String() string {
	return e.Email
}

type EmailInvalidError struct {
	Email string
}

func (e EmailInvalidError) Error() string {
	return fmt.Sprintf("%s is not a valid email address", e.Email)
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	if !emailRegex.MatchString(email) {
		return false
	}
	return true
}
