package resources

import "github.com/google/uuid"

type AuthUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}
