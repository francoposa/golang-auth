package db

import (
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
	"log"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgAuthZResourceRepo struct {
	db *sqlx.DB
}

func NewPGAuthZResourceTypeRepo(db *sqlx.DB) repos.ResourceRepo {
	return &pgAuthZResourceRepo{db: db}
}

var insertResourceStatement = `
INSERT INTO authz_resource_type (id, name, description) 
VALUES ($1, $2, $3)
RETURNING id, name, description
`

func (r *pgAuthZResourceRepo) Create(resource *resources.AuthZResourceType) (*resources.AuthZResourceType, error) {
	var id uuid.UUID
	var name string
	var description string
	err := r.db.QueryRowx(
		insertResourceStatement,
		resource.ID,
		resource.Name,
		resource.Description,
	).Scan(&id, &name, &description)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.ResourceNameAlreadyExistsError{ResourceName: resource.Name}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.AuthZResourceType{ID: id, Name: name, Description: description}, err
}

var selectResourceByNameStatement = `
SELECT * FROM authz_resource_type
WHERE name=$1
`

func (r *pgAuthZResourceRepo) Get(name string) (*resources.AuthZResourceType, error) {
	var id uuid.UUID
	var retrievedName string
	var description string
	err := r.db.QueryRowx(selectResourceByNameStatement,
		name,
	).Scan(&id, &retrievedName, &description)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &repos.ResourceNameNotFoundError{ResourceName: name}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.AuthZResourceType{ID: id, Name: retrievedName, Description: description}, err
}
