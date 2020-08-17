package repos

import (
	"fmt"
	"golang-auth/usecases/resources"
)

type ResourceRepo interface {
	Create(resource *resources.AuthZResourceType) (*resources.AuthZResourceType, error)
	Get(name string) (*resources.AuthZResourceType, error)
}

type ResourceNameNotFoundError struct {
	ResourceName string
}

func (e *ResourceNameNotFoundError) Error() string {
	return fmt.Sprintf("No Name found with name %s", e.ResourceName)
}

type ResourceNameAlreadyExistsError struct {
	ResourceName string
}

func (e *ResourceNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("Name already exists with name %s", e.ResourceName)
}
