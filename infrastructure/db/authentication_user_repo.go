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

type pgAuthNUserModel struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

func (model pgAuthNUserModel) toResource() *resources.AuthNUser {
	return &resources.AuthNUser{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
	}
}

type pgAuthNUserRepo struct {
	db     *sqlx.DB
	hasher usecases.Hasher
}

func NewPGAuthNUserRepo(db *sqlx.DB, hasher usecases.Hasher) *pgAuthNUserRepo {
	return &pgAuthNUserRepo{db: db, hasher: hasher}
}

var insertAuthNUserStatement = `
INSERT INTO authentication_user (id, username, email, password) 
VALUES ($1, $2, $3, $4)
RETURNING id, username, email
`

func (r *pgAuthNUserRepo) Create(user *resources.AuthNUser, password string) (*resources.AuthNUser, error) {
	hashedPassword, err := r.hasher.Hash(password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var au pgAuthNUserModel
	err = r.db.QueryRowx(
		insertAuthNUserStatement,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
	).StructScan(&au)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.DuplicateAuthNUserForUsernameError{Username: user.Username}
		}
		log.Print(err)
		return nil, err
	}
	return au.toResource(), err
}

var selectAuthNUserByUsernameStatement = `
SELECT * FROM authentication_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Get(username string) (*resources.AuthNUser, error) {
	var au pgAuthNUserModel
	err := r.db.QueryRowx(selectAuthNUserByUsernameStatement, username).StructScan(&au)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &repos.AuthNUserNotFoundForUsernameError{Username: username}
		}
		log.Print(err)
		return nil, err
	}
	return au.toResource(), nil
}

func (r *pgAuthNUserRepo) Verify(username string, password string) (bool, error) {
	var au pgAuthNUserModel
	err := r.db.QueryRowx(selectAuthNUserByUsernameStatement, username).StructScan(&au)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, &repos.AuthNUserNotFoundForUsernameError{Username: username}
		}
		log.Print(err)
		return false, err
	}

	return r.hasher.Verify(password, au.Password)
}
