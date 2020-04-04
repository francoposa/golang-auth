package interfaces

import "github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"

type ClientStore interface {
	Create() (resources.Client, error)
	Get(id string) (resources.Client, error)
}
