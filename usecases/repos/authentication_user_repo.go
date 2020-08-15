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

type AuthNUserUsernameNotFoundError struct {
	Username string
}

func (e AuthNUserUsernameNotFoundError) Error() string {
	return fmt.Sprintf("No AuthNUser found with username %s", e.Username)
}

type AuthNUsernameAlreadyExistsError struct {
	Username string
}

func (e AuthNUsernameAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNUser already exists with username %s", e.Username)
}
