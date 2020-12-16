package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pgTools "github.com/francoposa/go-tools/postgres"
	sqlxTools "github.com/francoposa/go-tools/postgres/sqlx"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"golang-auth/authentication/application/server"
	"golang-auth/authentication/infrastructure/crypto"
	"golang-auth/authentication/infrastructure/db"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		pgConfig := pgTools.Config{
			Host:                  viper.GetString("postgres.host"),
			Port:                  viper.GetInt("postgres.port"),
			Username:              viper.GetString("postgres.username"),
			Password:              viper.GetString("postgres.password"),
			Database:              viper.GetString("postgres.database"),
			ApplicationName:       viper.GetString("postgres.application"),
			ConnectTimeoutSeconds: viper.GetInt("postgres.connectTimeoutSeconds"),
			SSLMode:               viper.GetString("postgres.sslMode"),
		}
		sqlxDB := sqlxTools.MustConnect(pgConfig)

		hasher := crypto.NewDefaultArgon2PassHasher()
		userRepo := db.NewPGUserRepo(sqlxDB, hasher)
		userHandler := server.NewUserHandler(userRepo)
		loginHandler := server.NewLoginHandler(
			viper.GetString("ui_app.urls.login"),
		)

		router := chi.NewRouter()

		// Suggested basic middleware stack from chi's docs
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		// Routing to API handlers
		router.Route("/api/v1/login", func(router chi.Router) {
			router.Post("/", loginHandler.InitializeLogin)
		})

		router.Route("/api/v1/users", func(router chi.Router) {
			router.Post("/", userHandler.Create)
			router.Get("/{id}", userHandler.Get)
			//router.Post("/authenticate", userHandler.Authenticate)
		})

		handler := cors.Default().Handler(router)
		handler = handlers.LoggingHandler(os.Stdout, handler)

		host := viper.GetString("server.host")
		port := viper.GetString("server.port")
		readTimeout := viper.GetInt("server.timeout.read")
		writeTimeout := viper.GetInt("server.timeout.write")
		idleTimeout := viper.GetInt("server.timeout.idle")

		srv := &http.Server{
			Handler:      router,
			Addr:         host + ":" + port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		}

		fmt.Printf("running http server on port %s...\n", port)
		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// HTTP Server
	serverCmd.PersistentFlags().String("server.host", "", "")
	err := viper.BindPFlag(
		"server.host", serverCmd.PersistentFlags().Lookup("server.host"),
	)
	serverCmd.PersistentFlags().String("server.port", "", "")
	err = viper.BindPFlag(
		"server.port", serverCmd.PersistentFlags().Lookup("server.port"),
	)
	serverCmd.PersistentFlags().String("server.timeout.read", "", "")
	err = viper.BindPFlag(
		"server.timeout.read", serverCmd.PersistentFlags().Lookup("server.timeout.read"),
	)
	serverCmd.PersistentFlags().String("server.timeout.write", "", "")
	err = viper.BindPFlag(
		"server.timeout.write", serverCmd.PersistentFlags().Lookup("server.timeout.write"),
	)
	serverCmd.PersistentFlags().String("server.timeout.idle", "", "")
	err = viper.BindPFlag(
		"server.timeout.idle", serverCmd.PersistentFlags().Lookup("server.timeout.idle"),
	)

	// Postgres
	serverCmd.PersistentFlags().String("postgres.host", "", "")
	err = viper.BindPFlag(
		"postgres.host", serverCmd.PersistentFlags().Lookup("postgres.host"),
	)
	serverCmd.PersistentFlags().String("postgres.port", "", "")
	err = viper.BindPFlag(
		"postgres.port", serverCmd.PersistentFlags().Lookup("postgres.port"),
	)
	serverCmd.PersistentFlags().String("postgres.username", "", "")
	err = viper.BindPFlag(
		"postgres.username", serverCmd.PersistentFlags().Lookup("postgres.username"),
	)
	serverCmd.PersistentFlags().String("postgres.password", "", "")
	err = viper.BindPFlag(
		"postgres.password", serverCmd.PersistentFlags().Lookup("postgres.password"),
	)
	serverCmd.PersistentFlags().String("postgres.database", "", "")
	err = viper.BindPFlag(
		"postgres.database", serverCmd.PersistentFlags().Lookup("postgres.database"),
	)
	serverCmd.PersistentFlags().String("postgres.application", "", "")
	err = viper.BindPFlag(
		"postgres.application", serverCmd.PersistentFlags().Lookup("postgres.application"),
	)
	serverCmd.PersistentFlags().String("postgres.connectTimeoutSeconds", "", "")
	err = viper.BindPFlag(
		"postgres.connectTimeoutSeconds", serverCmd.PersistentFlags().Lookup("postgres.connectTimeoutSeconds"),
	)
	serverCmd.PersistentFlags().String("postgres.sslMode", "", "")
	err = viper.BindPFlag(
		"postgres.sslMode", serverCmd.PersistentFlags().Lookup("postgres.sslMode"),
	)

	// Login & Register UI
	serverCmd.PersistentFlags().String("ui_app.urls.login", "", "")
	err = viper.BindPFlag(
		"ui_app.urls.login", serverCmd.PersistentFlags().Lookup("ui_app.urls.login"),
	)
	serverCmd.PersistentFlags().String("ui_app.urls.register", "", "")
	err = viper.BindPFlag(
		"ui_app.urls.login", serverCmd.PersistentFlags().Lookup("ui_app.urls.register"),
	)

	if err != nil {
		panic(err)
	}
}
