package repos

import (
	"github.com/google/uuid"

	"golang-auth/usecases/resources"
)

type AuthNRoleRepo interface {
	GetByName(role string) (*resources.AuthNRole, error)
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

type DuplicateAuthNRole struct {
	errMsg string
}

func (e *DuplicateAuthNRole) Error() string {
	return e.errMsg
}
