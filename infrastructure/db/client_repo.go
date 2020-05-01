package db

import (
	"golang-auth/usecases/resources"
	"net/url"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type pgClientModel struct {
	ID          *uuid.UUID
	Secret      *uuid.UUID
	RedirectURI string `db:"redirect_uri"`
	Public      bool
}

func (model pgClientModel) toResource() *resources.Client {
	uri, _ := url.Parse(model.RedirectURI)
	return &resources.Client{
		ID:          model.ID,
		Secret:      model.Secret,
		RedirectURI: uri,
		Public:      model.Public,
	}
}

type pgClientRepo struct {
	db *sqlx.DB
}

func NewPGClientRepo(db *sqlx.DB) *pgClientRepo {
	return &pgClientRepo{db: db}
}

var insertClientStatement = `
INSERT INTO client (id, secret, redirect_uri, public) 
VALUES ($1, $2, $3, $4)
RETURNING id, secret, redirect_uri, public
`

func (r *pgClientRepo) Create(client *resources.Client) (*resources.Client, error) {
	var c pgClientModel
	err := r.db.QueryRowx(
		insertClientStatement,
		client.ID,
		client.Secret,
		client.RedirectURI.String(),
		client.Public,
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

func (r *pgClientRepo) Get(id *uuid.UUID) (*resources.Client, error) {
	var c pgClientModel
	err := r.db.QueryRowx(selectClientByIDStatement, id).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return c.toResource(), err
}
