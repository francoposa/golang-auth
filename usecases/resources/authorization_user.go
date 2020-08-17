package resources

import "github.com/google/uuid"

type AuthZUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func NewAuthZUser(username, email string) *AuthZUser {
	id := uuid.New()
	return &AuthZUser{ID: id, Username: username, Email: email}
}
