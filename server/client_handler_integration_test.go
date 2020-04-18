package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"golang-auth/db"
	"golang-auth/entities/resources"
)

func setupTestClientHandler(t *testing.T) *mux.Router {
	sqlxDB := db.SetUpDB(t)
	clientRepo, _ := db.SetUpClientRepo(t, sqlxDB)
	clientHandler := ClientHandler{repo: &clientRepo}

	router := mux.NewRouter()
	router.HandleFunc("/credentials/", clientHandler.CreateClient).Methods("POST")

	return router
}

func TestClientHandler_POSTClient(t *testing.T) {
	clientHandler := setupTestClientHandler(t)

	t.Run("POST new client - success", func(t *testing.T) {
		response := httptest.NewRecorder()
		clientHandler.ServeHTTP(response, newPOSTClientRequest(t))
		createdClient := resources.Client{}
		err := json.NewDecoder(response.Body).Decode(&createdClient)
		if err != nil {
			t.Errorf("Unable to parse response body from server %q into Client", response.Body)
		}
	})

	db.TearDownDB(t)
}

func newPOSTClientRequest(t *testing.T) *http.Request {
	t.Helper()
	body := map[string]string{"Domain": "example"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/credentials/"), bytes.NewBuffer(jsonBody))
	return req
}
