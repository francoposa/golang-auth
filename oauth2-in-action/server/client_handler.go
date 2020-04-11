package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/francojposa/golang-auth/oauth2-in-action/entities/interfaces"
)

type ClientHandler struct {
	repo interfaces.ClientRepo
}

func NewClientHandler(repo interfaces.ClientRepo) *ClientHandler {
	return &ClientHandler{repo: repo}
}

func (c *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	client, err := c.repo.Get(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(client)
	}

}
