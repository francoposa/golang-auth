package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"golang-auth/authentication/domain"
)

type PGAuthNUserRepo struct {
	db     *sqlx.DB
	hasher domain.Hasher
}

func NewPGAuthNUserRepo(db *sqlx.DB, hasher domain.Hasher) *PGAuthNUserRepo {
	return &PGAuthNUserRepo{db: db, hasher: hasher}
}

var insertAuthNUser = `
INSERT INTO authn_user (id, username, email, password) 
VALUES ($1, $2, $3, $4)
RETURNING id, username, email
`

func (r *PGAuthNUserRepo) Create(user *domain.AuthNUser, password string) (*domain.AuthNUser, error) {
	hashedPassword, err := r.hasher.Hash(password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var id uuid.UUID
	var username string
	var email string

	err = r.db.QueryRowx(
		insertAuthNUser,
		user.ID,
		user.Username,
		user.Email.String(),
		hashedPassword,
	).Scan(&id, &username, &email)

	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" {
				key, value := GetAlreadyExistsErrorKeyValue(pqError)
				return nil, domain.AuthNUserAlreadyExistsError{Field: key, Value: value}
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.AuthNUser{
		ID:       id,
		Username: username,
		Email:    domain.EmailAddress{Email: email},
	}, nil
}

const selectAuthNUserByUsername = `
SELECT id, username, email
FROM authn_user
WHERE username=$1
`

func (r *PGAuthNUserRepo) Get(username string) (*domain.AuthNUser, error) {
	var id uuid.UUID
	var email string

	err := r.db.QueryRow(
		selectAuthNUserByUsername,
		username,
	).Scan(&id, &username, &email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.AuthNUserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.AuthNUser{
		ID:       id,
		Username: username,
		Email:    domain.EmailAddress{Email: email},
	}, nil
}

const selectAuthNUserByUsernameWithPassword = `
SELECT id, username, email, password
FROM authn_user
WHERE username=$1
`

func (r *PGAuthNUserRepo) Verify(username string, password string) (bool, error) {
	var id uuid.UUID
	var email string
	var hashedPassword string

	err := r.db.QueryRowx(
		selectAuthNUserByUsernameWithPassword,
		username,
	).Scan(&id, &username, &email, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.AuthNUserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Println(err)
		return false, err
	}

	return r.hasher.Verify(password, hashedPassword)
}
