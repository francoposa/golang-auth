package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

type pgAuthNRoleRepo struct {
	db *sqlx.DB
}

func NewPGAuthNRoleRepo(db *sqlx.DB) repos.AuthNRoleRepo {
	return &pgAuthNRoleRepo{db: db}
}

var insertRoleStatement = `
INSERT INTO authentication_role (id, role)
VALUES ($1, $2)
RETURNING id, role
`

func (r *pgAuthNRoleRepo) Create(role *resources.AuthNRole) (*resources.AuthNRole, error) {
	var createdRole resources.AuthNRole
	err := r.db.QueryRowx(
		insertRoleStatement,
		role.ID,
		role.Role,
	).StructScan(&createdRole)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.DuplicateAuthNRole{Role: role.Role}
		}
		log.Print(err)
		return nil, err
	}
	return &createdRole, err
}

var selectRoleByRoleStatement = `
SELECT * FROM authentication_role
WHERE role=$1
`

func (r *pgAuthNRoleRepo) Get(role string) (*resources.AuthNRole, error) {
	var createdRole resources.AuthNRole
	err := r.db.QueryRowx(selectRoleByRoleStatement, role).StructScan(&createdRole)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &repos.AuthNRoleNotFoundError{Role: role}
		}
		log.Print(err)
		return nil, err
	}
	return &createdRole, err
}
