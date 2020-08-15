package repos

import (
	"fmt"
	"golang-auth/usecases/resources"
)

type ResourceRepo interface {
	Create(resource *resources.Resource) (*resources.Resource, error)
	Get(name string) (*resources.Resource, error)
}

type ResourceNameNotFoundError struct {
	ResourceName string
}

func (e *ResourceNameNotFoundError) Error() string {
	return fmt.Sprintf("No ResourceName found with name %s", e.ResourceName)
}

type ResourceNameAlreadyExistsError struct {
	ResourceName string
}

func (e *ResourceNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("ResourceName already exists with name %s", e.ResourceName)
}
