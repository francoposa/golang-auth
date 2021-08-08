package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	uuid "github.com/satori/go.uuid"

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

func (h *LoginHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login, err := h.repo.GetByID(id)
	var notFoundErr domain.LoginNotFoundError
	if errors.As(err, &notFoundErr) {
		w.WriteHeader(http.StatusNotFound)
		body := map[string]string{"errorMessage": err.Error()}
		json.NewEncoder(w).Encode(body)
		return
	}

	err = json.NewEncoder(w).Encode(
		HttpReadLogin{
			ID:          login.ID,
			RedirectURL: login.RedirectURL,
			Status:      login.Status,
			Attempts:    login.Attempts,
			CSRFToken:   login.CSRFToken,
			CreatedAt:   login.CreatedAt,
			UpdatedAt:   login.UpdatedAt,
		},
	)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *LoginHandler) InitializeLogin(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	login := domain.NewLogin(*r.URL, token)
	login, err := h.repo.Create(login)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	body := map[string]string{"login_id": login.ID.String()}
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// TODO finish implementing for redirect-based login apps
	//query := url.Values{}
	//query.Add("login-id", login.ID.String())
	//h.loginURL.RawQuery = query.Encode()
	//
	//http.Redirect(w, r, h.loginURL.String(), http.StatusFound)
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
