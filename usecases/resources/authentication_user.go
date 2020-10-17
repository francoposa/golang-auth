package resources

import (
	"github.com/google/uuid"
)

type AuthNUser struct {
	ID       uuid.UUID
	Username string
	Email    EmailAddress
}

func NewAuthNUser(username string, email EmailAddress) *AuthNUser {
	id := uuid.New()
	return &AuthNUser{ID: id, Username: username, Email: email}
}
