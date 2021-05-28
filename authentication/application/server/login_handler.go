package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"

	"golang-auth/authentication/domain"
)

type LoginHandler struct {
	repo     domain.LoginRepo
	userRepo domain.UserRepo
	loginURL url.URL
}

func NewLoginHandler(
	repo domain.LoginRepo,
	userRepo domain.UserRepo,
	loginURL url.URL,
) *LoginHandler {
	return &LoginHandler{repo, userRepo, loginURL}
}

func (h *LoginHandler) InitializeLogin(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	login := domain.NewLogin(*r.URL, token)
	login, err := h.repo.Create(login)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	query := url.Values{}
	query.Add("login-id", login.ID.String())
	h.loginURL.RawQuery = query.Encode()

	http.Redirect(w, r, h.loginURL.String(), http.StatusFound)
}

func (h *LoginHandler) VerifyLogin(w http.ResponseWriter, r *http.Request) {
	httpUser := HttpAuthenticateUser{}
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified, err := h.userRepo.VerifyPassword(httpUser.Username, httpUser.Password)
	if err != nil && errors.Is(
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
		body := map[string]string{"errorMessage": "Username or Password is incorrect"}
		json.NewEncoder(w).Encode(body)
		return
	}
	w.WriteHeader(http.StatusOK)
}
