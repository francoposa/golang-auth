package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"golang-auth/db"
	"golang-auth/server"
)

// authserverCmd represents the authserver command
var authserverCmd = &cobra.Command{
	Use:   "authserver",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		pgConfig := db.NewDefaultPostgresConfig("oauth2_in_action")
		sqlxDB := db.MustConnect(pgConfig)

		clientRepo := db.PGClientRepo{DB: sqlxDB}
		clientHandler := server.NewClientHandler(&clientRepo)

		router := mux.NewRouter()
		router.HandleFunc("/credentials/", clientHandler.CreateClient).Methods("POST")

		srv := &http.Server{
			Handler: router,
			Addr:    "127.0.0.1:5000",
			// Good practice: enforce timeouts for servers
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		fmt.Println("running http server on port 5000")
		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(authserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
