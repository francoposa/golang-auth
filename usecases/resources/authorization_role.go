package resources

import "github.com/google/uuid"

type AuthZRole struct {
	ID   uuid.UUID
	Name string
}

func NewAuthZRole(name string) *AuthZRole {
	id := uuid.New()
	return &AuthZRole{ID: id, Name: name}
}
