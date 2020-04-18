package server

import (
	"encoding/json"
	"net/http"

	"golang-auth/entities/interfaces"
	"golang-auth/entities/resources"
)

type httpPOSTClient struct {
	Domain string
}

type ClientHandler struct {
	repo interfaces.ClientRepo
}

func NewClientHandler(repo interfaces.ClientRepo) *ClientHandler {
	return &ClientHandler{repo: repo}
}

func (c *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	postedClient := httpPOSTClient{}
	err := json.NewDecoder(r.Body).Decode(&postedClient)
	if err != nil || postedClient.Domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := resources.NewClient(postedClient.Domain)
	repoClient, err := c.repo.Create(client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(repoClient)
	}
}

// Going to come back to this - obviously shouldn't have an
// unsecured method of listing client credentials

//func (c *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id := vars["id"]
//	uid, err := uuid.Parse(id)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	client, err := c.repo.Get(uid)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//	} else {
//		json.NewEncoder(w).Encode(client)
//	}
//
//}
