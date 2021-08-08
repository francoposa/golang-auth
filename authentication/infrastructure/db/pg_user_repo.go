package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"

	"golang-auth/authentication/domain"
)

type PGUserRepo struct {
	db     *sqlx.DB
	hasher domain.Hasher
}

func NewPGUserRepo(db *sqlx.DB, hasher domain.Hasher) *PGUserRepo {
	return &PGUserRepo{db: db, hasher: hasher}
}

const insertUser = `
INSERT INTO authn_user (
	id, username, email, password, enabled, created_at, updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, username, email, enabled, created_at, updated_at
`

func (r *PGUserRepo) Create(user *domain.User, password *domain.Password) (*domain.User, error) {

	hashedPassword, err := r.hasher.Hash(string(*password))
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	var username string
	var email string
	var enabled bool
	var created_at time.Time
	var updated_at time.Time

	err = r.db.QueryRow(
		insertUser,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
		user.Enabled,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id, &username, &email, &enabled, &created_at, &updated_at)

	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" {
				key, value := GetAlreadyExistsErrorKeyValue(pqError)
				return nil, domain.UserAlreadyExistsError{Field: key, Value: value}
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   enabled,
		CreatedAt: created_at.UTC(),
		UpdatedAt: updated_at.UTC(),
	}, nil
}

const selectUserByID = `
SELECT id, username, email, enabled, created_at, updated_at
FROM authn_user
WHERE id=$1
`

func (r *PGUserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var username string
	var email string
	var enabled bool
	var created_at time.Time
	var updated_at time.Time

	err := r.db.QueryRow(
		selectUserByID,
		id.String(),
	).Scan(&id, &username, &email, &enabled, &created_at, &updated_at)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.UserNotFoundError{
				Field: "id",
				Value: id.String(),
			}
		}
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   enabled,
		CreatedAt: created_at.UTC(),
		UpdatedAt: updated_at.UTC(),
	}, nil
}

const selectUserByUsername = `
SELECT id, username, email, enabled, created_at, updated_at
FROM authn_user
WHERE username=$1
`

func (r *PGUserRepo) GetByUsername(username string) (*domain.User, error) {
	var id uuid.UUID
	var email string
	var enabled bool
	var created_at time.Time
	var updated_at time.Time

	err := r.db.QueryRow(
		selectUserByUsername,
		username,
	).Scan(&id, &username, &email, &enabled, &created_at, &updated_at)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.UserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   enabled,
		CreatedAt: created_at.UTC(),
		UpdatedAt: updated_at.UTC(),
	}, nil
}

const selectUserByUsernameWithPassword = `
SELECT id, username, email, password
FROM authn_user
WHERE username=$1
`

func (r *PGUserRepo) VerifyPassword(username string, password string) (bool, error) {
	var id uuid.UUID
	var email string
	var hashedPassword string

	err := r.db.QueryRowx(
		selectUserByUsernameWithPassword,
		username,
	).Scan(&id, &username, &email, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.UserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		return false, err
	}

	return r.hasher.Verify(password, hashedPassword)
}
