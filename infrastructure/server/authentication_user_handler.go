package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang-auth/usecases/repos"
)

type httpUserAuthentication struct {
	Username string
	Password string
}

type AuthNUserHandler struct {
	repo repos.AuthNUserRepo
}

func NewAuthNUserHandler(repo repos.AuthNUserRepo) *AuthNUserHandler {
	return &AuthNUserHandler{repo: repo}
}

func (h *AuthNUserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	postedUserAuth := httpUserAuthentication{}
	err := json.NewDecoder(r.Body).Decode(&postedUserAuth)
	if err != nil || postedUserAuth.Username == "" || postedUserAuth.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := h.repo.Verify(postedUserAuth.Username, postedUserAuth.Password)
	// Handler AuthNUserNotFound
	if errors.Is(err, repos.AuthNUsernameNotFoundError{Username: postedUserAuth.Username}) {
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
	} else {
		w.WriteHeader(http.StatusOK)
	}

}
