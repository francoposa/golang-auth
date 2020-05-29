package resources

import "github.com/google/uuid"

type AuthNRole struct {
	ID   uuid.UUID
	Role string
}

func NewAuthNRole(role string) *AuthNRole {
	id := uuid.New()
	return &AuthNRole{ID: id, Role: role}
}
