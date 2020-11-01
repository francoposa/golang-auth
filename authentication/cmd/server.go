package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"golang-auth/authentication/application/server"
	"golang-auth/authentication/infrastructure/crypto"
	"golang-auth/authentication/infrastructure/db"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		templatePattern := filepath.Join(wd, "/application/web/templates/*")
		baseTemplatePath := filepath.Join(wd, "/application/web/templates/base.gohtml")

		templates := server.NewTemplates(templatePattern, baseTemplatePath)

		templateRenderer := server.NewTemplateRenderer(templates, "base")

		pgConfig := db.NewDefaultPostgresConfig("examplecom_auth")
		sqlxDB := db.MustConnect(pgConfig)

		hasher := crypto.NewDefaultArgon2PassHasher()

		authNUserRepo := db.NewPGAuthNUserRepo(sqlxDB, hasher)
		authNUserHandler := server.NewUserHandler(authNUserRepo)

		authNWebHandler := server.NewAuthNWebHandler(templateRenderer, "sign-in.gohtml", "sign-up.gohtml")

		router := mux.NewRouter()
		// Routing to FileServer handler for static web assets
		// Choose the folder to serve
		appRootDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		httpStaticAssetsDir := http.Dir(fmt.Sprintf("/%s/application/web/static/", appRootDir))
		staticRoute := "/static/"
		router.PathPrefix(staticRoute).Handler(http.StripPrefix(staticRoute, http.FileServer(httpStaticAssetsDir)))

		// Routing to HTML Template and API handlers
		router.HandleFunc("/login", authNWebHandler.GetLogin).Methods("GET")
		router.HandleFunc("/register", authNWebHandler.GetRegister).Methods("GET")
		router.HandleFunc("/login", authNUserHandler.Authenticate).Methods("POST")
		router.HandleFunc("/register", authNUserHandler.Create).Methods("POST")

		handler := cors.Default().Handler(router)
		handler = handlers.LoggingHandler(os.Stdout, handler)

		srv := &http.Server{
			Handler: handler,
			Addr:    "127.0.0.1:5001",
			// Good practice: enforce timeouts for servers
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		fmt.Println("running http server on port 5001")
		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
