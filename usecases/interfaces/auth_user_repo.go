package interfaces

import (
	"golang-auth/usecases/resources"
)

type AuthUserRepo interface {
	Get(username string) (*resources.AuthUser, error)
	Create(user *resources.AuthUser, password string) (*resources.AuthUser, error)
	Verify(username string, password string) (bool, error)
}
