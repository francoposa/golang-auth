package repos

import (
	"fmt"
	"golang-auth/usecases/resources"

	"github.com/google/uuid"
)

type AuthZUserRepo interface {
	Get(id uuid.UUID) (*resources.AuthZUser, error)
	Create(user *resources.AuthZUser) (*resources.AuthZUser, error)
}

type AuthZUserIDNotFoundError struct {
	ID string
}

func (e AuthZUserIDNotFoundError) Error() string {
	return fmt.Sprintf("No AuthZUser found with id %s", e.ID)
}

type AuthZUsernameAlreadyExistsError struct {
	Username string
}

func (e AuthZUsernameAlreadyExistsError) Error() string {
	return fmt.Sprintf("AuthNUser already exists with username %s", e.Username)
}
