package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang-auth/authentication-identity-user-mgmt/domain"
)

type UserHandler struct {
	repo domain.UserRepo
}

func NewUserHandler(repo domain.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

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

	user, err := domain.NewUser(httpUser.Username, httpUser.Email)
	var usernameErr domain.UsernameInvalidError
	var emailErr domain.EmailInvalidError
	if errors.As(err, &usernameErr) || errors.As(err, &emailErr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	createdUser, err := h.repo.Create(user, httpUser.Password)
	var passwordErr domain.PasswordInvalidError
	var existsErr domain.UserAlreadyExistsError
	if errors.As(err, &passwordErr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}
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

	w.WriteHeader(http.StatusOK)

}
