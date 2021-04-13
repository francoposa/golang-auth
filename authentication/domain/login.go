package domain

import (
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

const LoginStatusInitialized = "initialized"

type Login struct {
	ID          uuid.UUID
	RedirectURL url.URL
	Status      string
	Attempts    int
	CSRFToken   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewLogin(redirectURL url.URL, csrfToken string) *Login {
	id := uuid.NewV4()
	now := time.Now().UTC()
	return &Login{
		ID:          id,
		RedirectURL: redirectURL,
		Status:      LoginStatusInitialized,
		Attempts:    0,
		CSRFToken:   csrfToken,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
