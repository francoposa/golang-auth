package repos

import (
	"fmt"

	"github.com/google/uuid"

	"golang-auth/usecases/resources"
)

type AuthZRoleRepo interface {
	GetByName(name string) (*resources.AuthZRole, error)
	GetByID(id uuid.UUID) (*resources.AuthZRole, error)
	Create(role *resources.AuthZRole) (*resources.AuthZRole, error)
}

type AuthZRoleNameNotFoundError struct {
	Name string
}

func (e AuthZRoleNameNotFoundError) Error() string {
	return fmt.Sprintf("No AuthZRole found with name %s", e.Name)
}

type AuthZRoleNameAlreadyExistsError struct {
	Name string
}

func (e AuthZRoleNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNRole already exists with name %s", e.Name)
}
