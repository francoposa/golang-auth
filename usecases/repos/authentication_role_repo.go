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

type AuthNRoleNotFoundError struct {
	errMsg string
}

func NewAuthNRoleNotFoundError(errMsg string) AuthNRoleNotFoundError {
	return AuthNRoleNotFoundError{errMsg: errMsg}
}

func (e AuthNRoleNotFoundError) Error() string {
	return e.errMsg
}

type AuthNRoleNameAlreadyExistsError struct {
	RoleName string
}

func (e AuthNRoleNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNRole already exists with rolename %s", e.RoleName)
}
