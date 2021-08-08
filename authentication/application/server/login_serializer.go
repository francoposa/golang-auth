package server

import (
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

type HttpReadLogin struct {
	ID          uuid.UUID `json:"id"`
	RedirectURL url.URL   `json:"redirect_url"`
	Status      string    `json:"status"`
	Attempts    int       `json:"attempts"`
	CSRFToken   string    `json:"csrf_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
