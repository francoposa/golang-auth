package db

import (
	"log"

	"github.com/google/uuid"

	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgAuthZRoleRepo struct {
	db *sqlx.DB
}

func NewPGAuthZRoleRepo(db *sqlx.DB) repos.AuthZRoleRepo {
	return &pgAuthZRoleRepo{db: db}
}

var insertRoleStatement = `
INSERT INTO authz_role (id, name)
VALUES ($1, $2)
RETURNING id, name
`

func (r *pgAuthZRoleRepo) Create(role *resources.AuthZRole) (*resources.AuthZRole, error) {
	var id uuid.UUID
	var name string
	err := r.db.QueryRowx(
		insertRoleStatement,
		role.ID,
		role.Name,
	).Scan(&id, &name)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, repos.AuthZRoleNameAlreadyExistsError{Name: name}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.AuthZRole{ID: id, Name: name}, err
}

var selectRoleByIDStatement = `
SELECT * FROM authz_role
WHERE id=$1
`

func (r *pgAuthZRoleRepo) GetByID(id uuid.UUID) (*resources.AuthZRole, error) {
	var retrievedID uuid.UUID
	var name string
	err := r.db.QueryRowx(
		selectRoleByIDStatement,
		id,
	).Scan(&retrievedID, &name)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, repos.AuthZRoleNameNotFoundError{Name: name}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.AuthZRole{ID: retrievedID, Name: name}, err
}

var selectRoleByNameStatement = `
SELECT * FROM authz_role
WHERE name=$1
`

func (r *pgAuthZRoleRepo) GetByName(name string) (*resources.AuthZRole, error) {
	var id uuid.UUID
	var retrievedName string
	err := r.db.QueryRowx(
		selectRoleByNameStatement,
		name,
	).Scan(&id, &retrievedName)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, repos.AuthZRoleNameNotFoundError{Name: name}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.AuthZRole{ID: id, Name: retrievedName}, err
}
