package resources

import "github.com/google/uuid"

type AuthNUser struct {
	ID       uuid.UUID
	Username string
	Email    string
	Role     AuthNRole
}

func NewAuthNUser(username, email string, role AuthNRole) *AuthNUser {
	id := uuid.New()
	return &AuthNUser{ID: id, Username: username, Email: email, Role: role}
}
