package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"golang-auth/infrastructure/db"
	"golang-auth/infrastructure/server"
)

// authzserverCmd represents the authzserver command
var authzserverCmd = &cobra.Command{
	Use:   "authzserver",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {



		pgConfig := db.NewDefaultPostgresConfig("examplecom_auth")
		sqlxDB := db.MustConnect(pgConfig)



		clientRepo := db.NewPGClientRepo(sqlxDB)
		clientHandler := server.NewClientHandler(clientRepo)

		authHandler := server.AuthorizationHandler{}

		router := mux.NewRouter()
		router.HandleFunc("/authorize", authHandler.Authorize).Methods("GET", "POST")

		router.HandleFunc("/client", clientHandler.Create).Methods("POST")

		handler := cors.Default().Handler(router)

		srv := &http.Server{
			Handler: handler,
			Addr:    "127.0.0.1:5002",
			// Good practice: enforce timeouts for servers
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		fmt.Println("running http server on port 5002")
		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(authzserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authzserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authzserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
