package resources

import (
	"github.com/google/uuid"
)

// Resource is an abstract entity in the ExampleCom system which may be operated upon
// OAuth Clients request authorization to operate on Resources on behalf of a Resource Owner.
type Resource struct {
	ID           uuid.UUID
	ResourceName string
}

func NewResource(resourceName string) *Resource {
	return &Resource{
		ID:           uuid.New(),
		ResourceName: resourceName,
	}
}
