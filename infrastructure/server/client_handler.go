package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
)

type httpPOSTRequestClient struct {
	RedirectURI string
	Public      bool
	FirstParty  bool
}

type httpResponseClient struct {
	ID          uuid.UUID
	Secret      *uuid.UUID
	RedirectURI string
	Public      bool
	FirstParty  bool
}

func responseClientFromResource(resource *resources.Client) httpResponseClient {
	return httpResponseClient{
		ID:          resource.ID,
		Secret:      resource.Secret,
		RedirectURI: resource.RedirectURI.String(),
		Public:      resource.Public,
		FirstParty:  resource.FirstParty,
	}
}

type ClientHandler struct {
	repo repos.ClientRepo
}

func NewClientHandler(repo repos.ClientRepo) *ClientHandler {
	return &ClientHandler{repo: repo}
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	postedClient := httpPOSTRequestClient{}
	err := json.NewDecoder(r.Body).Decode(&postedClient)
	if err != nil || postedClient.RedirectURI == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client, err := resources.NewClient(postedClient.RedirectURI, postedClient.Public, postedClient.FirstParty)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repoClient, err := h.repo.Create(client)
	responseClient := responseClientFromResource(repoClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responseClient)
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
//	client, err := c.repo.GetByName(uid)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//	} else {
//		json.NewEncoder(w).Encode(client)
//	}
//
//}
