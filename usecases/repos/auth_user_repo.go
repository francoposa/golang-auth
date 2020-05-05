package repos

import (
	"fmt"

	"golang-auth/usecases/resources"
)

type AuthUserRepo interface {
	Get(username string) (*resources.AuthUser, error)
	Create(user *resources.AuthUser, password string) (*resources.AuthUser, error)
	Verify(username string, password string) (bool, error)
}

type AuthUserNotFoundForUsernameError struct {
	Username string
}

func (e *AuthUserNotFoundForUsernameError) Error() string {
	return fmt.Sprintf("No AuthUser found for username %s", e.Username)
}

type DuplicateAuthUserForUsernameError struct {
	Username string
}

func (e *DuplicateAuthUserForUsernameError) Error() string {
	return fmt.Sprintf("AuthUser already exists for username %s", e.Username)
}
