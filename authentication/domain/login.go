package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Login struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
