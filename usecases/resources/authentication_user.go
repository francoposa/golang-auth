package resources

import "github.com/google/uuid"

type AuthNUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func NewAuthNUser(username, email string) *AuthNUser {
	id := uuid.New()
	return &AuthNUser{ID: id, Username: username, Email: email}
}
