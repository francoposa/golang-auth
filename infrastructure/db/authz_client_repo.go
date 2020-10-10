package db

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

type pgClientModel struct {
	ID          uuid.UUID
	Secret      *uuid.UUID
	RedirectURI string `db:"redirect_uri"`
	Public      bool
	FirstParty  bool `db:"first_party"`
}

func (model pgClientModel) toResource() *resources.AuthZClient {
	uri, _ := url.Parse(model.RedirectURI)
	return &resources.AuthZClient{
		ID:          model.ID,
		Secret:      model.Secret,
		RedirectURI: uri,
		Public:      model.Public,
		FirstParty:  model.FirstParty,
	}
}

type pgAuthZClientRepo struct {
	db *sqlx.DB
}

func NewPGAuthZClientRepo(db *sqlx.DB) repos.AuthZClientRepo {
	return &pgAuthZClientRepo{db: db}
}

var insertClientStatement = `
INSERT INTO authz_client (id, secret, redirect_uri, public, first_party) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, secret, redirect_uri, public, first_party
`

func (r *pgAuthZClientRepo) Create(client *resources.AuthZClient) (*resources.AuthZClient, error) {
	var c pgClientModel
	err := r.db.QueryRowx(
		insertClientStatement,
		client.ID,
		client.Secret,
		client.RedirectURI.String(),
		client.Public,
		client.FirstParty,
	).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return c.toResource(), err
}

var selectClientByIDStatement = `
SELECT * FROM authz_client
WHERE id=$1
`

func (r *pgAuthZClientRepo) Get(id uuid.UUID) (*resources.AuthZClient, error) {
	var c pgClientModel
	err := r.db.QueryRowx(selectClientByIDStatement, id).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return c.toResource(), err
}
