package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
			return nil, repos.NewAuthNRoleAlreadyExistsError("AuthNRole already exists")
		}
		log.Print(err)
		return nil, err
	}
	return &createdRole, err
}

var selectRoleByIDStatement = `
SELECT * FROM authentication_role
WHERE id=$1
`

func (r *pgAuthNRoleRepo) GetByID(id uuid.UUID) (*resources.AuthNRole, error) {
	var retreivedRole resources.AuthNRole
	err := r.db.QueryRowx(selectRoleByIDStatement, id).StructScan(&retreivedRole)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			errMsg := fmt.Sprintf("No AuthNRole found with id %v", id)
			return nil, repos.NewAuthNRoleNotFoundError(errMsg)
		}
		log.Print(err)
		return nil, err
	}
	return &retreivedRole, err
}

var selectRoleByRoleStatement = `
SELECT * FROM authentication_role
WHERE role=$1
`

func (r *pgAuthNRoleRepo) GetByName(roleName string) (*resources.AuthNRole, error) {
	var retreivedRole resources.AuthNRole
	err := r.db.QueryRowx(selectRoleByRoleStatement, roleName).StructScan(&retreivedRole)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			errMsg := fmt.Sprintf("No AuthNRole found with role name %s", roleName)
			return nil, repos.NewAuthNRoleNotFoundError(errMsg)
		}
		log.Print(err)
		return nil, err
	}
	return &retreivedRole, err
}
