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

func setupTestAuthNUserHandler(t *testing.T, sqlxDB *sqlx.DB) *mux.Router {
	t.Helper()

	authNRoleRepo, _ := db.SetUpAuthNRoleRepo(t, sqlxDB)
	authNUserRepo, _ := db.SetUpAuthNUserRepo(t, sqlxDB, authNRoleRepo)
	authNUserHandler := AuthNUserHandler{repo: authNUserRepo}

	router := mux.NewRouter()
	router.HandleFunc("/login", authNUserHandler.Authenticate).Methods("POST")

	return router
}

func TestAuthNUserHandler_Authenticate(t *testing.T) {
	assertions := assert.New(t)
	sqlxDB, closeDB := db.SetUpDB(t)
	defer closeDB(t, sqlxDB)

	AuthNUserHandler := setupTestAuthNUserHandler(t, sqlxDB)

	t.Run("HTTP 200 for correct username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_pass"}
		AuthNUserHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(200, response.Code, "Expected HTTP 200, got: %d", response.Code)
	})

	t.Run("HTTP 401 for incorrect username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_badpass"}
		AuthNUserHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(401, response.Code, "Expected HTTP 401, got: %d", response.Code)

		body = map[string]string{"username": "domtoretto_badusername", "password": "domtoretto_pass"}
		AuthNUserHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(401, response.Code, "Expected HTTP 401, got: %d", response.Code)
	})

}

func newPOSTUserAuthenticateRequest(t *testing.T, body map[string]string) *http.Request {
	t.Helper()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/login"), bytes.NewBuffer(jsonBody))
	return req
}
