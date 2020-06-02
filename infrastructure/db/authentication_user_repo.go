package db

import (
	"fmt"
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

func NewPGAuthNUserRepo(db *sqlx.DB, hasher usecases.Hasher, authNRoleRepo repos.AuthNRoleRepo) repos.AuthNUserRepo {
	return &pgAuthNUserRepo{db: db, hasher: hasher, authNRoleRepo: authNRoleRepo}
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
			return nil, repos.NewAuthNUserAlreadyExistsError("AuthNUser already exists")
		}
		log.Print(err)
		return nil, err
	}

	role, err := r.authNRoleRepo.GetByID(role_id)
	if err != nil {
		return nil, err
	}

	return &resources.AuthNUser{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
	}, nil
}

var selectAuthNUserByUsernameStatement = `
SELECT id, username, email, role_id
FROM authentication_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Get(username string) (*resources.AuthNUser, error) {
	var id uuid.UUID
	var email string
	var role_id uuid.UUID

	err := r.db.QueryRowx(
		selectAuthNUserByUsernameStatement,
		username,
	).Scan(&id, &username, &email, &role_id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			errMsg := fmt.Sprintf("No AuthNUser found with username %s", username)
			return nil, repos.NewAuthNUserNotFoundError(errMsg)
		}
		log.Print(err)
		return nil, err
	}

	role, err := r.authNRoleRepo.GetByID(role_id)
	if err != nil {
		return nil, err
	}

	return &resources.AuthNUser{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
	}, nil
}

var selectAuthNUserByUsernameForPasswordVerificationStatement = `
SELECT id, username, email, password, role_id
FROM authentication_user
WHERE username=$1
`

func (r *pgAuthNUserRepo) Verify(username string, password string) (bool, error) {
	var id uuid.UUID
	var email string
	var hashedPassword string
	var role_id uuid.UUID

	err := r.db.QueryRowx(
		selectAuthNUserByUsernameForPasswordVerificationStatement,
		username,
	).Scan(&id, &username, &email, &hashedPassword, &role_id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			errMsg := fmt.Sprintf("No AuthNUser found with username %s", username)
			return false, repos.NewAuthNUserNotFoundError(errMsg)
		}
		log.Print(err)
		return false, err
	}

	return r.hasher.Verify(password, hashedPassword)
}
