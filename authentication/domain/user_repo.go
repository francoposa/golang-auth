package domain

import (
	"fmt"
)

type AuthNUserRepo interface {
	Get(username string) (*AuthNUser, error)
	Create(user *AuthNUser, password string) (*AuthNUser, error)
	Verify(username string, password string) (bool, error)
}

type AuthNUserNotFoundError struct {
	Field string
	Value string
}

func (e AuthNUserNotFoundError) Error() string {
	return fmt.Sprintf("No AuthNUser found with exists with %s=%s", e.Field, e.Value)
}

type AuthNUserAlreadyExistsError struct {
	Field string
	Value string
}

func (e AuthNUserAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNUser already exists with %s=%s", e.Field, e.Value)
}
