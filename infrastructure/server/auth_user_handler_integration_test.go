package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"

	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"golang-auth/infrastructure/db"
)

func setupTestAuthUserHandler(t *testing.T) *mux.Router {
	sqlxDB := db.SetUpDB(t)
	authUserRepo, _ := db.SetUpAuthUserRepo(t, sqlxDB)
	authUserHandler := AuthUserHandler{repo: authUserRepo}

	router := mux.NewRouter()
	router.HandleFunc("/user/authenticate/", authUserHandler.Authenticate).Methods("POST")

	return router
}

func TestAuthUserHandler_Authenticate(t *testing.T) {
	assertions := assert.New(t)
	authUserHandler := setupTestAuthUserHandler(t)

	t.Run("HTTP 200 for correct username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_pass"}
		authUserHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(200, response.Code, "Expected HTTP 200, got: %d", response.Code)
	})

	t.Run("HTTP 401 for incorrect username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_badpass"}
		authUserHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(401, response.Code, "Expected HTTP 401, got: %d", response.Code)

		body = map[string]string{"username": "domtoretto_badusername", "password": "domtoretto_pass"}
		authUserHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(401, response.Code, "Expected HTTP 401, got: %d", response.Code)
	})

}

func newPOSTUserAuthenticateRequest(t *testing.T, body map[string]string) *http.Request {
	t.Helper()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/authenticate/"), bytes.NewBuffer(jsonBody))
	return req
}
