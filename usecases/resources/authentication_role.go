package resources

import "github.com/google/uuid"

type AuthNRole struct {
	ID       uuid.UUID
	RoleName string
}

func NewAuthNRole(roleName string) *AuthNRole {
	id := uuid.New()
	return &AuthNRole{ID: id, RoleName: roleName}
}
