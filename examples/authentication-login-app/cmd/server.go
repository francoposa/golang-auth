package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/francoposa/golang-auth/examples/authentication-login-app/application/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		templatePattern := filepath.Join(wd, "/application/web/templates/*")
		baseTemplatePath := filepath.Join(wd, "/application/web/templates/base.gohtml")
		templates := server.NewTemplates(templatePattern, baseTemplatePath)
		templateRenderer := server.NewTemplateRenderer(templates)

		webHandler := server.NewWebHandler(templateRenderer)

		httpStaticAssetsDir := http.Dir(fmt.Sprintf("%s/application/web/static/", wd))
		staticRoute := "/static/"
		staticAssetHandler := http.StripPrefix(
			staticRoute,
			http.FileServer(httpStaticAssetsDir),
		)

		router := chi.NewRouter()

		// Suggested basic middleware stack from chi's docs
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		router.Handle(staticRoute+"*", staticAssetHandler)
		router.Get("/login", webHandler.GetLogin)
		router.Get("/register", webHandler.GetRegister)

		host := viper.GetString("server.host")
		port := viper.GetString("server.port")
		readTimeout := viper.GetInt("server.timeout.read")
		writeTimeout := viper.GetInt("server.timeout.write")

		srv := &http.Server{
			Handler:      router,
			Addr:         host + ":" + port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		}

		fmt.Printf("running http server on port %s...\n", port)
		log.Fatal(srv.ListenAndServe())

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
