package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"

	"golang-auth/authentication/domain"
)

type UserHandler struct {
	repo domain.UserRepo
}

func NewUserHandler(repo domain.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	retrievedUser, err := h.repo.GetByID(id)
	var notFoundErr domain.UserNotFoundError
	if errors.As(err, &notFoundErr) {
		w.WriteHeader(http.StatusNotFound)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	err = json.NewEncoder(w).Encode(
		HttpReadUser{
			ID:       retrievedUser.ID,
			Username: retrievedUser.Username,
			Email:    retrievedUser.Email,
		},
	)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	httpUser := HttpCreateUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = httpUser.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	user, password, err := domain.NewUser(httpUser.Username, httpUser.Email, httpUser.Password)
	var usernameErr domain.UsernameInvalidError
	var emailErr domain.EmailInvalidError
	var passwordErr domain.PasswordInvalidLengthError
	if errors.As(err, &usernameErr) || errors.As(err, &emailErr) || errors.As(err, &passwordErr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	createdUser, err := h.repo.Create(user, password)
	var existsErr domain.UserAlreadyExistsError
	if errors.As(err, &existsErr) {
		w.WriteHeader(http.StatusConflict)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	err = json.NewEncoder(w).Encode(
		HttpReadUser{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			Email:    createdUser.Email,
		},
	)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
