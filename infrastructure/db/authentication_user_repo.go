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
	db            *sqlx.DB
	hasher        usecases.Hasher
	authNRoleRepo repos.AuthNRoleRepo
}

func NewPGAuthNUserRepo(db *sqlx.DB, hasher usecases.Hasher) repos.AuthNUserRepo {
	return &pgAuthNUserRepo{db: db, hasher: hasher}
}

var insertAuthNUserStatement = `
INSERT INTO authentication_user (id, username, email, password, role_id) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, email, role_id
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
	var role_id uuid.UUID

	err = r.db.QueryRowx(
		insertAuthNUserStatement,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
		user.Role.ID,
	).Scan(&id, &username, &email, &role_id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.DuplicateAuthNUserForUsernameError{Username: user.Username}
		}
		log.Print(err)
		return nil, err
	}

	role, err := r.authNRoleRepo.GetByName()

	return &resources.AuthNUser{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
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
