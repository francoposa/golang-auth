package interfaces

import (
	"golang-auth/usecases/resources"

	"github.com/google/uuid"
)

type ClientRepo interface {
	Create(client *resources.Client) (*resources.Client, error)
	Get(id *uuid.UUID) (*resources.Client, error)
}
