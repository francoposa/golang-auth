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

type pgAuthUserModel struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

func (model pgAuthUserModel) toResource() *resources.AuthUser {
	return &resources.AuthUser{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
	}
}

type pgAuthUserRepo struct {
	db     *sqlx.DB
	hasher usecases.Hasher
}

func NewPGAuthUserRepo(db *sqlx.DB, hasher usecases.Hasher) *pgAuthUserRepo {
	return &pgAuthUserRepo{db: db, hasher: hasher}
}

var insertAuthuserStatement = `
INSERT INTO auth_user (id, username, email, password) 
VALUES ($1, $2, $3, $4)
RETURNING id, username, email
`

func (r *pgAuthUserRepo) Create(user *resources.AuthUser, password string) (*resources.AuthUser, error) {
	hashedPassword, err := r.hasher.Hash(password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var au pgAuthUserModel
	err = r.db.QueryRowx(
		insertAuthuserStatement,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
	).StructScan(&au)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.DuplicateAuthUserForUsernameError{Username: user.Username}
		}
		log.Print(err)
		return nil, err
	}
	return au.toResource(), err
}

var selectAuthUserByUsernameStatement = `
SELECT * FROM auth_user WHERE username=$1
`

func (r *pgAuthUserRepo) Get(username string) (*resources.AuthUser, error) {
	var au pgAuthUserModel
	err := r.db.QueryRowx(selectAuthUserByUsernameStatement, username).StructScan(&au)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &repos.AuthUserNotFoundForUsernameError{Username: username}
		}
		log.Print(err)
		return nil, err
	}
	return au.toResource(), nil
}

func (r *pgAuthUserRepo) Verify(username string, password string) (bool, error) {
	var au pgAuthUserModel
	err := r.db.QueryRowx(selectAuthUserByUsernameStatement, username).StructScan(&au)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, &repos.AuthUserNotFoundForUsernameError{Username: username}
		}
		log.Print(err)
		return false, err
	}

	return r.hasher.Verify(password, au.Password)
}
