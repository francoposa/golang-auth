package db

import (
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
	"log"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgResourceRepo struct {
	db *sqlx.DB
}

func NewPGResourceRepo(db *sqlx.DB) repos.ResourceRepo {
	return &pgResourceRepo{db: db}
}

var insertResourceStatement = `
INSERT INTO resource (id, resource_name) 
VALUES ($1, $2)
RETURNING id, resource_name
`

func (r *pgResourceRepo) Create(resource *resources.Resource) (*resources.Resource, error) {
	var id uuid.UUID
	var resourceName string
	err := r.db.QueryRowx(
		insertResourceStatement,
		resource.ID,
		resource.ResourceName,
	).Scan(&id, &resourceName)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.ResourceNameAlreadyExistsError{ResourceName: resource.ResourceName}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.Resource{ID: id, ResourceName: resourceName}, err
}

var selectResourceByNameStatement = `
SELECT * FROM resource
WHERE resource_name=$1
`

func (r *pgResourceRepo) Get(resourceName string) (*resources.Resource, error) {
	var id uuid.UUID
	var retrievedResourceName string
	err := r.db.QueryRowx(selectResourceByNameStatement,
		resourceName,
	).Scan(&id, &retrievedResourceName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &repos.ResourceNameNotFoundError{ResourceName: resourceName}
		}
		log.Print(err)
		return nil, err
	}
	return &resources.Resource{ID: id, ResourceName: retrievedResourceName}, err
}
