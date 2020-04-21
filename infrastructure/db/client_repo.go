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

type PGClientRepo struct {
	db *sqlx.DB
}

func NewPGClientRepo(db *sqlx.DB) *PGClientRepo {
	return &PGClientRepo{db: db}
}

var insertStatement = `
INSERT INTO client (id, secret, domain) 
VALUES ($1, $2, $3)
RETURNING id, secret, domain`

func (r *PGClientRepo) Create(client *resources.Client) (*resources.Client, error) {
	var cm pgClientModel
	err := r.db.QueryRowx(
		insertStatement,
		client.ID,
		client.Secret,
		client.Domain,
	).StructScan(&cm)
	if err != nil {
		return nil, err
	}
	return cm.toEntity(), err
}

var selectByIDStatement = `
SELECT * FROM client
WHERE id=$1
`

func (r *PGClientRepo) Get(id uuid.UUID) (*resources.Client, error) {
	var cm pgClientModel
	err := r.db.QueryRowx(selectByIDStatement, id).StructScan(&cm)
	if err != nil {
		return nil, err
	}
	return cm.toEntity(), err
}

func (model pgClientModel) toEntity() *resources.Client {
	return &resources.Client{
		ID:     model.ID,
		Secret: model.Secret,
		Domain: model.Domain,
	}
}
