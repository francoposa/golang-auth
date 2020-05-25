package repos

import (
	"fmt"

	"golang-auth/usecases/resources"
)

type AuthNRoleRepo interface {
	Get(role string) (*resources.AuthNRole, error)
	Create(role *resources.AuthNRole) (*resources.AuthNRole, error)
}

type AuthNRoleNotFoundError struct {
	Role string
}

func (e *AuthNRoleNotFoundError) Error() string {
	return fmt.Sprintf("No AuthNRole found with role %s", e.Role)
}

type DuplicateAuthNRole struct {
	Role string
}

func (e *DuplicateAuthNRole) Error() string {
	return fmt.Sprintf("AuthNRole already exists with role %s", e.Role)
}
