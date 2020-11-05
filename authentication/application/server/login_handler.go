package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang-auth/authentication/domain"
)

type LoginHandler struct {
	userRepo domain.UserRepo
}

func NewLoginHandler(userRepo domain.UserRepo) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	httpUser := HttpAuthenticateUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := h.userRepo.VerifyPassword(httpUser.Username, httpUser.Password)
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
