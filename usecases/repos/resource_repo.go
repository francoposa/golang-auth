package repos

import (
	"fmt"
	"golang-auth/usecases/resources"
)

type ResourceRepo interface {
	Create(client *resources.Resource) (*resources.Resource, error)
	Get(name string) (*resources.Resource, error)
}

type ResourceNotFoundForNameError struct {
	Name string
}

func (e *ResourceNotFoundForNameError) Error() string {
	return fmt.Sprintf("No Resource found with name %s", e.Name)
}

type DuplicateResourceForNameError struct {
	Name string
}

func (e *DuplicateResourceForNameError) Error() string {
	return fmt.Sprintf("Resource already exists with name %s", e.Name)
}
