package db

import (
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

var insertAuthNUserStatement = `
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
		insertAuthNUserStatement,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
	).Scan(&id, &username, &email)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, repos.AuthNUsernameAlreadyExistsError{Username: username}
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

var selectAuthNUserByUsernameStatement = `
SELECT id, username, email
FROM authn_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Get(username string) (*resources.AuthNUser, error) {
	var id uuid.UUID
	var email string

	err := r.db.QueryRowx(
		selectAuthNUserByUsernameStatement,
		username,
	).Scan(&id, &username, &email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, repos.AuthNUserUsernameNotFoundError{username}
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

var selectAuthNUserByUsernameForPasswordVerificationStatement = `
SELECT id, username, email, password
FROM authn_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Verify(username string, password string) (bool, error) {
	var id uuid.UUID
	var email string
	var hashedPassword string

	err := r.db.QueryRowx(
		selectAuthNUserByUsernameForPasswordVerificationStatement,
		username,
	).Scan(&id, &username, &email, &hashedPassword)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, repos.AuthNUserUsernameNotFoundError{username}
		}
		log.Print(err)
		return false, err
	}

	return r.hasher.Verify(password, hashedPassword)
}
