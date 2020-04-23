package db

import (
	"golang-auth/usecases/resources"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type pgClientModel struct {
	ID     uuid.UUID
	Secret uuid.UUID
	Domain string
}

func (model pgClientModel) toResource() *resources.Client {
	return &resources.Client{
		ID:     model.ID,
		Secret: model.Secret,
		Domain: model.Domain,
	}
}

type pgClientRepo struct {
	db *sqlx.DB
}

func NewPGClientRepo(db *sqlx.DB) *pgClientRepo {
	return &pgClientRepo{db: db}
}

var insertClientStatement = `
INSERT INTO client (id, secret, domain) 
VALUES ($1, $2, $3)
RETURNING id, secret, domain
`

func (r *pgClientRepo) Create(client *resources.Client) (*resources.Client, error) {
	var c pgClientModel
	err := r.db.QueryRowx(
		insertClientStatement,
		client.ID,
		client.Secret,
		client.Domain,
	).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return c.toResource(), err
}

var selectClientByIDStatement = `
SELECT * FROM client
WHERE id=$1
`

func (r *pgClientRepo) Get(id uuid.UUID) (*resources.Client, error) {
	var c pgClientModel
	err := r.db.QueryRowx(selectClientByIDStatement, id).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return c.toResource(), err
}
