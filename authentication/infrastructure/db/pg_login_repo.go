package db

import (
	"log"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"golang-auth/authentication/domain"
)

type PGLoginRepo struct {
	db       *sqlx.DB
	userRepo domain.UserRepo
}

func NewPGLoginRepo(db *sqlx.DB, userRepo domain.UserRepo) *PGLoginRepo {
	return &PGLoginRepo{
		db:       db,
		userRepo: userRepo,
	}
}

const insertLogin = `
INSERT INTO authn_login (
	id, redirect_url, status, attempts, csrf_token, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, redirect_url, status, attempts, csrf_token, created_at, updated_at`

func (r *PGLoginRepo) Create(login *domain.Login) (*domain.Login, error) {
	var id uuid.UUID
	var rawRedirectURL string
	var status string
	var attempts int
	var csrfToken string
	var createdAt time.Time
	var UpdatedAt time.Time

	err := r.db.QueryRow(
		insertLogin,
		login.ID,
		login.RedirectURL.String(),
		login.Status,
		login.Attempts,
		login.CSRFToken,
		login.CreatedAt,
		login.UpdatedAt,
	).Scan(&id, &rawRedirectURL, &status, &attempts, &csrfToken, &createdAt, &UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	redirectURL, err := url.Parse(rawRedirectURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &domain.Login{
		ID:          id,
		RedirectURL: *redirectURL,
		Status:      status,
		Attempts:    attempts,
		CSRFToken:   csrfToken,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   UpdatedAt.UTC(),
	}, nil

}
