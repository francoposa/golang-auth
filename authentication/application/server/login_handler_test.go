package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	pgTools "github.com/francoposa/go-tools/postgres"
	sqlTools "github.com/francoposa/go-tools/postgres/database_sql"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/francoposa/golang-auth/authentication/infrastructure/db"
)

func setupTestUserHandler(t *testing.T, sqlxDB *sqlx.DB) chi.Router {
	t.Helper()

	userRepo, _ := db.SetUpUserRepo(t, sqlxDB)
	userHandler := UserHandler{repo: userRepo}

	router := chi.NewRouter()
	router.Post("/api/v1/users/authenticate", userHandler.Authenticate)

	return router
}

func TestAuthNUserHandler_Authenticate(t *testing.T) {
	assertions := assert.New(t)
	superUserPGConfig := pgTools.Config{
		Host:                  "localhost",
		Port:                  5432,
		Username:              "postgres",
		Password:              "",
		Database:              "postgres",
		ApplicationName:       "",
		ConnectTimeoutSeconds: 5,
		SSLMode:               "disable",
	}
	dbName := pgTools.RandomDBName("auth_test")

	sqlDB, err := db.SetUpDB(t, dbName, superUserPGConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlTools.TearDownDB(t, sqlDB, superUserPGConfig)

	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	userHandler := setupTestUserHandler(t, sqlxDB)

	t.Run("HTTP 200 for correct username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_password12345"}
		userHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(200, response.Code)
	})

	t.Run("HTTP 401 for incorrect username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_badpassword"}
		userHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(401, response.Code)

		body = map[string]string{"username": "domtoretto_badusername", "password": "domtoretto_password12345"}
		userHandler.ServeHTTP(response, newPOSTUserAuthenticateRequest(t, body))

		assertions.Equal(401, response.Code)
	})

}

func newPOSTUserAuthenticateRequest(t *testing.T, body map[string]string) *http.Request {
	t.Helper()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/api/v1/users/authenticate"),
		bytes.NewBuffer(jsonBody),
	)
	return req
}
