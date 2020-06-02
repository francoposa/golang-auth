package repos

import (
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

func NewAuthNRoleNotFoundError(errMsg string) *AuthNRoleNotFoundError {
	return &AuthNRoleNotFoundError{errMsg: errMsg}
}

func (e *AuthNRoleNotFoundError) Error() string {
	return e.errMsg
}

type AuthNRoleAlreadyExistsError struct {
	errMsg string
}

func NewAuthNRoleAlreadyExistsError(errMsg string) *AuthNRoleAlreadyExistsError {
	return &AuthNRoleAlreadyExistsError{errMsg: errMsg}
}

func (e *AuthNRoleAlreadyExistsError) Error() string {
	return e.errMsg
}
