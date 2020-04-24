package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang-auth/usecases/interfaces"
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
	db         *sqlx.DB
	passHasher interfaces.PassHasher
}

var insertAuthuserStatement = `
INSERT INTO auth_user (id, username, email, password) 
VALUES ($1, $2, $3, $4)
RETURNING id, username, email
`

func (r *pgAuthUserRepo) Create(user *resources.AuthUser, password string) (*resources.AuthUser, error) {
	hashedPassword, err := r.passHasher.Hash(password)
	if err != nil {
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
		return nil, err
	}
	return au.toResource(), nil
}

func (r *pgAuthUserRepo) Verify(username string, password string) (bool, error) {
	var au pgAuthUserModel
	err := r.db.QueryRowx(selectAuthUserByUsernameStatement, username).StructScan(&au)
	if err != nil {
		return false, err
	}

	return r.passHasher.Verify(password, au.Password)
}
