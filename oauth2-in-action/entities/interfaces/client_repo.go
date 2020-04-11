package interfaces

import (
	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
	"github.com/google/uuid"
)

type ClientRepo interface {
	Create(client *resources.Client) (*resources.Client, error)
	Get(id uuid.UUID) (*resources.Client, error)
}
