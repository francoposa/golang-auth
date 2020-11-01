package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang-auth/authentication/domain"
)

type AuthNUserHandler struct {
	repo domain.UserRepo
}

func NewUserHandler(repo domain.UserRepo) *AuthNUserHandler {
	return &AuthNUserHandler{repo: repo}
}

func (h *AuthNUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	httpUser := HttpCreateUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)

	if err != nil || !httpUser.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := domain.NewUser(httpUser.Username, httpUser.Email)
	var usernameErr domain.UsernameInvalidError
	var emailErr domain.EmailInvalidError
	if errors.As(err, &usernameErr) || errors.As(err, &emailErr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body := map[string]string{"error_message": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	createdUser, err := h.repo.Create(user, httpUser.Password)
	var passwordErr domain.PasswordInvalidError
	var existsErr domain.UserAlreadyExistsError
	if errors.As(err, &passwordErr) || errors.As(err, &existsErr) {
		w.WriteHeader(http.StatusConflict)
		body := map[string]string{"error_message": err.Error()}
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

func (h *AuthNUserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	httpUser := HttpAuthenticateUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := h.repo.Verify(httpUser.Username, httpUser.Password)
	if errors.Is(
		err,
		domain.UserNotFoundError{Field: "username", Value: httpUser.Username}) {
		verified = false
	} else if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !verified {
		w.WriteHeader(http.StatusUnauthorized)
		body := map[string]string{"error_message": "Username or Password is incorrect"}
		json.NewEncoder(w).Encode(body)
		return
	}
	w.WriteHeader(http.StatusOK)
}
