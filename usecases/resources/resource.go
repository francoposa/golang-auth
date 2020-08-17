package resources

import (
	"github.com/google/uuid"
)

// AuthZResourceType is an abstract entity in the ExampleCom system which may be operated upon
// OAuth Clients request authorization to operate on Resources on behalf of a AuthZResourceType Owner.
type AuthZResourceType struct {
	ID          uuid.UUID
	Name        string
	Description string
}

func NewAuthZResourceType(name, description string) *AuthZResourceType {
	return &AuthZResourceType{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
	}
}
