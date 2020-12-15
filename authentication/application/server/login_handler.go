package server

import "net/http"

type LoginHandler struct {
	LoginURL string
}

func NewLoginHandler(loginURL string) *LoginHandler {
	return &LoginHandler{LoginURL: loginURL}
}

func (h *LoginHandler) InitializeLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.LoginURL, http.StatusFound)
}
