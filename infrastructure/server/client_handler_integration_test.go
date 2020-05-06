package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"golang-auth/infrastructure/db"
)

func setupTestClientHandler(t *testing.T, sqlxDB *sqlx.DB) *mux.Router {
	t.Helper()
	clientRepo, _ := db.SetUpClientRepo(t, sqlxDB)
	clientHandler := ClientHandler{repo: clientRepo}

	router := mux.NewRouter()
	router.HandleFunc("/client/", clientHandler.Create).Methods("POST")

	return router

}

func TestClientHandler_POSTClient(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := db.SetUpDB(t)
	defer closeDB(t, sqlxDB)
	clientHandler := setupTestClientHandler(t, sqlxDB)

	t.Run("POST new client - success", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := httpPOSTRequestClient{RedirectURI: "telnyx.com", Public: true}
		clientHandler.ServeHTTP(response, newPOSTClientRequest(t, body))

		assertions.Equal(201, response.Code, "Expected HTTP 201, got: %d", response.Code)

		createdClient := httpResponseClient{}
		err := json.NewDecoder(response.Body).Decode(&createdClient)
		if err != nil {
			t.Errorf("Unable to parse response body from server %q into Client", response.Body)
		}
	})
}

func newPOSTClientRequest(t *testing.T, body httpPOSTRequestClient) *http.Request {
	t.Helper()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/client/"), bytes.NewBuffer(jsonBody))
	return req
}
