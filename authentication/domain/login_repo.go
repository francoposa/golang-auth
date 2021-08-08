package domain

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type LoginRepo interface {
	GetByID(id uuid.UUID) (*Login, error)
	Create(login *Login) (*Login, error)
}

type LoginNotFoundError struct {
	Field string
	Value string
}

func (e LoginNotFoundError) Error() string {
	return fmt.Sprintf("No Login found with %s=%s", e.Field, e.Value)
}
