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
}

func (model pgAuthUserModel) toResource() *resources.AuthUser {
	return &resources.AuthUser{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
	}
}

type PGAuthUserRepo struct {
	db         *sqlx.DB
	passHasher interfaces.PassHasher
}

func (p PGAuthUserRepo) Get(id uuid.UUID) (*resources.AuthUser, error) {
	panic("implement me")
}

var insertAuthuserStatement = `
INSERT INTO auth_user (id, username, email, password) 
VALUES ($1, $2, $3, $4)
RETURNING id, username, email
`

func (p PGAuthUserRepo) Create(user *resources.AuthUser, password string) (*resources.AuthUser, error) {
	hashedPassword, err := p.passHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	var au pgAuthUserModel
	err = p.db.QueryRowx(
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

func (p PGAuthUserRepo) Verify(id uuid.UUID, password string) (bool, error) {
	panic("implement me")
}
