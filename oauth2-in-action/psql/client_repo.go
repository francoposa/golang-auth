package psql

import (
	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
	"github.com/jmoiron/sqlx"
)

type pgClientModel struct {
	ID     string
	Secret string
	Domain string
}

type PGClientRepo struct {
	Db *sqlx.DB
}

func (r *PGClientRepo) GetClient(ID string) (*resources.Client, error) {
	cm := pgClientModel{}
	err := r.Db.Get(&cm, "SELECT * FROM client WHERE id=$1", ID)
	if err != nil {
		return nil, err
	}
	return modelToResource(cm), err
}

func modelToResource(model pgClientModel) *resources.Client {
	return &resources.Client{
		ID:     model.ID,
		Secret: model.Secret,
		Domain: model.Domain,
	}
}
