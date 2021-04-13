package domain

//import uuid "github.com/satori/go.uuid"

type LoginRepo interface {
	Create(login *Login) (*Login, error)
	//GetByID(id uuid.UUID) (*User, error)
}
