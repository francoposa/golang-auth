package repos

import (
	"fmt"
	"golang-auth/usecases/resources"
)

type AuthNUserRepo interface {
	Get(username string) (*resources.AuthNUser, error)
	Create(user *resources.AuthNUser, password string) (*resources.AuthNUser, error)
	Verify(username string, password string) (bool, error)
}

type AuthNUsernameNotFoundError struct {
	Username string
}

func (e AuthNUsernameNotFoundError) Error() string {
	return fmt.Sprintf("No AuthNUser found with username %s", e.Username)
}

type AuthNUserAlreadyExistsError struct {
	Field string
	Value string
}

func (e AuthNUserAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNUser already exists with %s=%s", e.Field, e.Value)
}
