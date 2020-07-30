/*
Copyright © 2020 Franco Posa <francojposa@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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

	"golang-auth/infrastructure/crypto"
	"golang-auth/infrastructure/db"
	"golang-auth/infrastructure/server"
)

// authnserverCmd represents the authnserver command
var authnserverCmd = &cobra.Command{
	Use:   "authnserver",
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

		templatePattern := filepath.Join(wd, "/infrastructure/web/templates/*")
		baseTemplatePath := filepath.Join(wd, "infrastructure/web/templates/base.gohtml")

		templates := server.NewTemplates(templatePattern, baseTemplatePath)

		templateRenderer := server.NewTemplateRenderer(templates, "base")

		pgConfig := db.NewDefaultPostgresConfig("examplecom_auth")
		sqlxDB := db.MustConnect(pgConfig)

		hasher := crypto.NewDefaultArgon2PassHasher()

		authNRoleRepo := db.NewPGAuthNRoleRepo(sqlxDB)
		//role, _ := authNRoleRepo.Create(resources.NewAuthNRole("user"))

		authNUserRepo := db.NewPGAuthNUserRepo(sqlxDB, hasher, authNRoleRepo)
		//_, _ = authNUserRepo.Create(resources.NewAuthNUser("test", "test@test.com", role), "test")
		authNUserHandler := server.NewAuthNUserHandler(authNUserRepo)

		loginHandler := server.NewLoginHandler(authNUserRepo, templateRenderer, "login.gohtml")

		router := mux.NewRouter()
		// Routing to FileServer handler for static web assets
		// Choose the folder to serve
		appRootDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		httpStaticAssetsDir := http.Dir(fmt.Sprintf("/%s/infrastructure/web/static/", appRootDir))
		staticRoute := "/static/"
		router.PathPrefix(staticRoute).Handler(http.StripPrefix(staticRoute, http.FileServer(httpStaticAssetsDir)))

		// Routing to HTML Template and API handlers
		router.HandleFunc("/login", loginHandler.Get).Methods("GET")
		router.HandleFunc("/login", authNUserHandler.Authenticate).Methods("POST")

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
	rootCmd.AddCommand(authnserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authnserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authnserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
