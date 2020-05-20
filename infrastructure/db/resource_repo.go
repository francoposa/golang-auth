package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
	"log"
)

type pgResourceModel struct {
	ID   uuid.UUID
	Name string
}

func (model pgResourceModel) toResource() *resources.Resource {
	return &resources.Resource{
		ID:   model.ID,
		Name: model.Name,
	}
}

type pgResourceRepo struct {
	db *sqlx.DB
}

func NewPGResourceRepo(db *sqlx.DB) *pgResourceRepo {
	return &pgResourceRepo{db: db}
}

var insertResourceStatement = `
INSERT INTO resource (id, name) 
VALUES ($1, $2)
RETURNING id, name
`

func (r *pgResourceRepo) Create(resource *resources.Resource) (*resources.Resource, error) {
	var model pgResourceModel
	err := r.db.QueryRowx(
		insertResourceStatement,
		resource.ID,
		resource.Name,
	).StructScan(&model)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, &repos.DuplicateResourceForNameError{Name: resource.Name}
		}
		log.Print(err)
		return nil, err
	}
	return model.toResource(), err
}

var selectResourceByNameStatement = `
SELECT * FROM resource
WHERE name=$1
`

func (r *pgResourceRepo) Get(name string) (*resources.Resource, error) {
	var model pgResourceModel
	err := r.db.QueryRowx(selectResourceByNameStatement, name).StructScan(&model)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &repos.ResourceNotFoundForNameError{Name: name}
		}
		log.Print(err)
		return nil, err
	}
	return model.toResource(), err
}
