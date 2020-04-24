package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"golang-auth/infrastructure/db"
	"golang-auth/usecases/resources"
)

func setupTestClientHandler(t *testing.T) *mux.Router {
	sqlxDB := db.SetUpDB(t)
	clientRepo, _ := db.SetUpClientRepo(t, sqlxDB)
	clientHandler := ClientHandler{repo: clientRepo}

	router := mux.NewRouter()
	router.HandleFunc("/client/", clientHandler.Create).Methods("POST")

	return router
}

func TestClientHandler_POSTClient(t *testing.T) {
	assertions := assert.New(t)
	clientHandler := setupTestClientHandler(t)

	t.Run("POST new client - success", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"Domain": "example"}
		clientHandler.ServeHTTP(response, newPOSTClientRequest(t, body))

		assertions.Equal(201, response.Code, "Expected HTTP 201, got: %d", response.Code)

		createdClient := resources.Client{}
		err := json.NewDecoder(response.Body).Decode(&createdClient)
		if err != nil {
			t.Errorf("Unable to parse response body from server %q into Client", response.Body)
		}
	})
}

func newPOSTClientRequest(t *testing.T, body map[string]string) *http.Request {
	t.Helper()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/client/"), bytes.NewBuffer(jsonBody))
	return req
}
