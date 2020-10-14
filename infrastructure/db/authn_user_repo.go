package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"golang-auth/usecases"
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

type pgAuthNUserRepo struct {
	db     *sqlx.DB
	hasher usecases.Hasher
}

func NewPGAuthNUserRepo(db *sqlx.DB, hasher usecases.Hasher) repos.AuthNUserRepo {
	return &pgAuthNUserRepo{db: db, hasher: hasher}
}

var insertAuthNUser = `
INSERT INTO authn_user (id, username, email, password) 
VALUES ($1, $2, $3, $4)
RETURNING id, username, email
`

func (r *pgAuthNUserRepo) Create(user *resources.AuthNUser, password string) (*resources.AuthNUser, error) {
	hashedPassword, err := r.hasher.Hash(password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var id uuid.UUID
	var username string
	var email string

	err = r.db.QueryRowx(
		insertAuthNUser,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
	).Scan(&id, &username, &email)

	var pqError *pq.Error
	if errors.As(err, &pqError) {
		if pqError.Code == "23505" {
			key, value := GetAlreadyExistsErrorKeyValue(pqError)
			return nil, repos.AuthNUserAlreadyExistsError{Field: key, Value: value}
		}
		log.Print(err)
		return nil, err
	}

	return &resources.AuthNUser{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

const selectAuthNUserByUsername = `
SELECT id, username, email
FROM authn_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Get(username string) (*resources.AuthNUser, error) {
	var id uuid.UUID
	var email string

	err := r.db.QueryRow(
		selectAuthNUserByUsername,
		username,
	).Scan(&id, &username, &email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repos.AuthNUserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Print(err)
		return nil, err
	}

	return &resources.AuthNUser{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

const selectAuthNUserByUsernameWithPassword = `
SELECT id, username, email, password
FROM authn_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Verify(username string, password string) (bool, error) {
	var id uuid.UUID
	var email string
	var hashedPassword string

	err := r.db.QueryRowx(
		selectAuthNUserByUsernameWithPassword,
		username,
	).Scan(&id, &username, &email, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repos.AuthNUserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Print(err)
		return false, err
	}

	return r.hasher.Verify(password, hashedPassword)
}
