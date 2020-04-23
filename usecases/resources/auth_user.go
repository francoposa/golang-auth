package resources

import "github.com/google/uuid"

type AuthUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func NewAuthUser(username, email string) *AuthUser {
	id := uuid.New()
	return &AuthUser{ID: id, Username: username, Email: email}
}
