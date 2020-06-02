package repos

import (
	"golang-auth/usecases/resources"
)

type AuthNUserRepo interface {
	Get(username string) (*resources.AuthNUser, error)
	Create(user *resources.AuthNUser, password string) (*resources.AuthNUser, error)
	Verify(username string, password string) (bool, error)
}


type AuthNUserNotFoundError struct {
	errMsg string
}

func NewAuthNUserNotFoundError(errMsg string) *AuthNUserNotFoundError {
	return &AuthNUserNotFoundError{errMsg: errMsg}
}

func (e *AuthNUserNotFoundError) Error() string {
	return e.errMsg
}

type AuthNUserAlreadyExistsError struct {
	errMsg string
}

func NewAuthNUserAlreadyExistsError(errMsg string) *AuthNUserAlreadyExistsError {
	return &AuthNUserAlreadyExistsError{errMsg: errMsg}
}

func (e *AuthNUserAlreadyExistsError) Error() string {
	return e.errMsg
}
