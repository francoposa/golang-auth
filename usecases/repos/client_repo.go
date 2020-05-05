package repos

import (
	"github.com/google/uuid"

	"golang-auth/usecases/resources"
)

type ClientRepo interface {
	Create(client *resources.Client) (*resources.Client, error)
	Get(id *uuid.UUID) (*resources.Client, error)
}
