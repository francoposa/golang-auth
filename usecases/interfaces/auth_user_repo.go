package interfaces

import (
	"github.com/google/uuid"

	"golang-auth/usecases/resources"
)

type AuthUserRepo interface {
	Get(id uuid.UUID) (*resources.AuthUser, error)
	Create(user *resources.AuthUser, password string) (*resources.AuthUser, error)
	Verify(id uuid.UUID, password string) (bool, error)
}
