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
	id := chi.URLParam(r, "id")
	uid, err := uuid.FromString(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	retrievedUser, err := h.repo.GetByID(uid)
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

	user, err := domain.NewUser(httpUser.Username, httpUser.Email)
	var usernameErr domain.UsernameInvalidError
	var emailErr domain.EmailInvalidError
	if errors.As(err, &usernameErr) || errors.As(err, &emailErr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}
	var passwordErr domain.PasswordInvalidError
	err = domain.ValidatePasswordRequirements(httpUser.Password)
	if errors.As(err, &passwordErr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	createdUser, err := h.repo.Create(user, httpUser.Password)
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

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	httpUser := HttpAuthenticateUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := h.repo.VerifyPassword(httpUser.Username, httpUser.Password)
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
