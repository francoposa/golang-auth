package server

import (
	"errors"

	"github.com/google/uuid"
)

type HttpCreateUser struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (u *HttpCreateUser) Validate() error {
	if u.Password != u.ConfirmPassword {
		return errors.New("Passwords do not match")
	}
	return nil
}

type HttpReadUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type HttpAuthenticateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
