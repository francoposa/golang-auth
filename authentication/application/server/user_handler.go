package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang-auth/authentication/domain"
)

type httpAuthNUser struct {
	Username        string `json:"username"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
}

type AuthNUserHandler struct {
	repo domain.AuthNUserRepo
}

func NewAuthNUserHandler(repo domain.AuthNUserRepo) *AuthNUserHandler {
	return &AuthNUserHandler{repo: repo}
}

func (h *AuthNUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	httpUser := httpAuthNUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil || httpUser.Username == "" || httpUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *AuthNUserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	httpUser := httpAuthNUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil || httpUser.Username == "" || httpUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := h.repo.Verify(httpUser.Username, httpUser.Password)
	if errors.Is(
		err,
		domain.AuthNUserNotFoundError{Field: "username", Value: httpUser.Username}) {
		verified = false
	} else if err != nil {
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
