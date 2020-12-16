package db

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"

	"golang-auth/authentication/domain"
)

type PGAuthNUserRepo struct {
	db     *sqlx.DB
	hasher domain.Hasher
}

func NewPGUserRepo(db *sqlx.DB, hasher domain.Hasher) *PGAuthNUserRepo {
	return &PGAuthNUserRepo{db: db, hasher: hasher}
}

const insertUser = `
INSERT INTO authn_user (
	id, username, email, password, enabled, created_at, updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, username, email, enabled, created_at, updated_at
`

func (r *PGAuthNUserRepo) Create(user *domain.User, password string) (*domain.User, error) {
	if !validator.StringLength(
		password,
		strconv.Itoa(domain.MinPasswordLen),
		strconv.Itoa(domain.MaxPasswordLen),
	) {
		return nil, domain.PasswordInvalidError{}
	}

	hashedPassword, err := r.hasher.Hash(password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var id uuid.UUID
	var username string
	var email string
	var enabled bool
	var created_at time.Time
	var updated_at time.Time

	err = r.db.QueryRow(
		insertUser,
		user.ID,
		user.Username,
		user.Email,
		hashedPassword,
		user.Enabled,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id, &username, &email, &enabled, &created_at, &updated_at)

	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" {
				key, value := GetAlreadyExistsErrorKeyValue(pqError)
				return nil, domain.UserAlreadyExistsError{Field: key, Value: value}
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   enabled,
		CreatedAt: created_at.UTC(),
		UpdatedAt: updated_at.UTC(),
	}, nil
}

const selectUserByID = `
SELECT id, username, email, enabled, created_at, updated_at
FROM authn_user
WHERE id=$1
`

func (r *PGAuthNUserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var username string
	var email string
	var enabled bool
	var created_at time.Time
	var updated_at time.Time

	err := r.db.QueryRow(
		selectUserByID,
		id.String(),
	).Scan(&id, &username, &email, &enabled, &created_at, &updated_at)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.UserNotFoundError{
				Field: "id",
				Value: id.String(),
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   enabled,
		CreatedAt: created_at.UTC(),
		UpdatedAt: updated_at.UTC(),
	}, nil
}

const selectUserByUsername = `
SELECT id, username, email, enabled, created_at, updated_at
FROM authn_user
WHERE username=$1
`

func (r *PGAuthNUserRepo) GetByUsername(username string) (*domain.User, error) {
	var id uuid.UUID
	var email string
	var enabled bool
	var created_at time.Time
	var updated_at time.Time

	err := r.db.QueryRow(
		selectUserByUsername,
		username,
	).Scan(&id, &username, &email, &enabled, &created_at, &updated_at)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.UserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Println(err)
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Enabled:   enabled,
		CreatedAt: created_at.UTC(),
		UpdatedAt: updated_at.UTC(),
	}, nil
}

const selectUserByUsernameWithPassword = `
SELECT id, username, email, password
FROM authn_user
WHERE username=$1
`

func (r *PGAuthNUserRepo) VerifyPassword(username string, password string) (bool, error) {
	var id uuid.UUID
	var email string
	var hashedPassword string

	err := r.db.QueryRowx(
		selectUserByUsernameWithPassword,
		username,
	).Scan(&id, &username, &email, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.UserNotFoundError{
				Field: "username",
				Value: username,
			}
		}
		log.Println(err)
		return false, err
	}

	return r.hasher.Verify(password, hashedPassword)
}
