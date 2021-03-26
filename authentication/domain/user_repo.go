package domain

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type UserRepo interface {
	Create(user *User, password *Password) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	VerifyPassword(username string, password string) (bool, error)
}

type UserNotFoundError struct {
	Field string
	Value string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("No User found with %s=%s", e.Field, e.Value)
}

type UserAlreadyExistsError struct {
	Field string
	Value string
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("User already exists with %s=%s", e.Field, e.Value)
}
