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

type AuthNUserNotFoundForUsernameError struct {
	Username string
}

func (e *AuthNUserNotFoundForUsernameError) Error() string {
	return fmt.Sprintf("No AuthNUser found with username %s", e.Username)
}

type DuplicateAuthNUserForUsernameError struct {
	Username string
}

func (e *DuplicateAuthNUserForUsernameError) Error() string {
	return fmt.Sprintf("AuthNUser already exists with username %s", e.Username)
}
