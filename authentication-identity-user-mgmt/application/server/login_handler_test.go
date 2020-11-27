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
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"golang-auth/authentication-identity-user-mgmt/infrastructure/db"
)

func setupTestLoginHandler(t *testing.T, sqlxDB *sqlx.DB) *mux.Router {
	t.Helper()

	userRepo, _ := db.SetUpUserRepo(t, sqlxDB)
	loginHandler := LoginHandler{userRepo: userRepo}

	router := mux.NewRouter()
	router.HandleFunc("/login", loginHandler.Login).Methods("POST")

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
	dbName := sqlTools.RandomDBName("auth_test")

	sqlDB, err := db.SetUpDB(t, dbName, superUserPGConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlTools.TearDownDB(t, sqlDB, superUserPGConfig)

	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	loginHandler := setupTestLoginHandler(t, sqlxDB)

	t.Run("HTTP 200 for correct username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_password12345"}
		loginHandler.ServeHTTP(response, newPOSTLoginRequest(t, body))

		assertions.Equal(200, response.Code)
	})

	t.Run("HTTP 401 for incorrect username and password", func(t *testing.T) {
		response := httptest.NewRecorder()
		body := map[string]string{"username": "domtoretto", "password": "domtoretto_badpass"}
		loginHandler.ServeHTTP(response, newPOSTLoginRequest(t, body))

		assertions.Equal(401, response.Code)

		body = map[string]string{"username": "domtoretto_badusername", "password": "domtoretto_password12345"}
		loginHandler.ServeHTTP(response, newPOSTLoginRequest(t, body))

		assertions.Equal(401, response.Code)
	})

}

func newPOSTLoginRequest(t *testing.T, body map[string]string) *http.Request {
	t.Helper()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/login"),
		bytes.NewBuffer(jsonBody),
	)
	return req
}
