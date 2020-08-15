package repos

import (
	"fmt"

	"github.com/google/uuid"

	"golang-auth/usecases/resources"
)

type AuthNRoleRepo interface {
	GetByName(roleName string) (*resources.AuthNRole, error)
	GetByID(id uuid.UUID) (*resources.AuthNRole, error)
	Create(role *resources.AuthNRole) (*resources.AuthNRole, error)
}

type AuthNRoleNameNotFoundError struct {
	RoleName string
}

func (e AuthNRoleNameNotFoundError) Error() string {
	return fmt.Sprintf("No AuthNUser found with role name %s", e.RoleName)
}

type AuthNRoleNameAlreadyExistsError struct {
	RoleName string
}

func (e AuthNRoleNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNRole already exists with role name %s", e.RoleName)
}
