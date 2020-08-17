package repos

import (
	"github.com/google/uuid"

	"golang-auth/usecases/resources"
)

type AuthZClientRepo interface {
	Create(client *resources.AuthZClient) (*resources.AuthZClient, error)
	Get(id uuid.UUID) (*resources.AuthZClient, error)
}
