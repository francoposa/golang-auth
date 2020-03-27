package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/francojposa/golang-auth/auth0-jwt-example/server"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	jwtMiddleWare := server.NewJWTMiddleware()

	// Setup API handling on router
	// /health - health check
	router.Handle("/health", http.HandlerFunc(server.HealthHandler)).Methods("GET")
	// /get-token - get JWT
	router.Handle("/get-token", http.HandlerFunc(server.GetTokenHandler)).Methods("GET")

	// /locations - retrieve a list of We We locations a user can leave feedback on
	router.Handle("/locations",
		jwtMiddleWare.Handler(http.HandlerFunc(server.ListLocationsHandler)),
	).Methods("GET")
	// /locations/{slug}/feedback - which will capture user feedback on locations
	router.Handle("/locations/{slug}/feedback",
		jwtMiddleWare.Handler(http.HandlerFunc(server.AddLocationFeedback)),
	).Methods("POST")

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, router),
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("running http server on port 3000")
	log.Fatal(srv.ListenAndServe())

}