package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	pgTools "github.com/francoposa/go-tools/postgres"
	sqlxTools "github.com/francoposa/go-tools/postgres/sqlx"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
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

		csrfSecure := viper.GetBool(serverCSRFSecureFlag)
		csrfKey := []byte(viper.GetString(serverCSRFKeyFlag))

		// Routing to API handlers
		router.Route("/api/v1/login", func(router chi.Router) {
			//router.Post("/", loginHandler.InitializeLogin)
			router.With(csrf.Protect(
				csrfKey, csrf.Secure(csrfSecure), csrf.Path("/"),
			)).Get("/", loginHandler.InitializeLogin)
			router.With(csrf.Protect(
				csrfKey, csrf.Secure(csrfSecure), csrf.Path("/"),
			)).Put("/", loginHandler.VerifyLogin)
		})

		router.Route("/api/v1/users", func(router chi.Router) {
			router.Post("/", userHandler.Create)
			router.Get("/{id}", userHandler.Get)
			//router.Post("/authenticate", userHandler.Authenticate)
		})

		//corsRouter := cors.Default().Handler(router)
		corsRouter := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut},
			AllowCredentials: true,
		}).Handler(router)
		//handler = handlers.LoggingHandler(os.Stdout, handler)

		host := viper.GetString(serverHostFlag)
		port := viper.GetString(serverPortFlag)
		readTimeout := viper.GetInt(serverTimeoutReadFlag)
		writeTimeout := viper.GetInt(serverTimeoutWriteFlag)
		idleTimeout := viper.GetInt(serverTimeoutIdleFlag)

		srv := &http.Server{
			Handler:      corsRouter,
			Addr:         host + ":" + port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		}

		fmt.Printf("running http server on port %s...\n", port)
		log.Fatal(srv.ListenAndServe())
	},
}

const serverHostFlag = "server.host"
const serverPortFlag = "server.port"
const serverTimeoutReadFlag = "server.timeout.read"
const serverTimeoutWriteFlag = "server.timeout.write"
const serverTimeoutIdleFlag = "server.timeout.idle"
const serverCSRFSecureFlag = "server.csrf.secure"
const serverCSRFKeyFlag = "server.csrf.key"

func init() {
	rootCmd.AddCommand(serverCmd)

	// HTTP Server
	serverCmd.PersistentFlags().String(serverHostFlag, "", "")
	err := viper.BindPFlag(
		serverHostFlag, serverCmd.PersistentFlags().Lookup(serverHostFlag),
	)
	serverCmd.PersistentFlags().String(serverPortFlag, "", "")
	err = viper.BindPFlag(
		serverPortFlag, serverCmd.PersistentFlags().Lookup(serverPortFlag),
	)
	serverCmd.PersistentFlags().String(serverTimeoutReadFlag, "", "")
	err = viper.BindPFlag(
		serverTimeoutReadFlag, serverCmd.PersistentFlags().Lookup(serverTimeoutReadFlag),
	)
	serverCmd.PersistentFlags().String(serverTimeoutWriteFlag, "", "")
	err = viper.BindPFlag(
		serverTimeoutWriteFlag, serverCmd.PersistentFlags().Lookup(serverTimeoutWriteFlag),
	)
	serverCmd.PersistentFlags().String(serverTimeoutIdleFlag, "", "")
	err = viper.BindPFlag(
		serverTimeoutIdleFlag, serverCmd.PersistentFlags().Lookup(serverTimeoutIdleFlag),
	)
	// HTTP Server CSRF
	serverCmd.PersistentFlags().String(serverCSRFSecureFlag, "", "")
	err = viper.BindPFlag(
		serverCSRFSecureFlag, serverCmd.PersistentFlags().Lookup(serverCSRFSecureFlag),
	)
	serverCmd.PersistentFlags().String(serverCSRFKeyFlag, "", "")
	err = viper.BindPFlag(
		serverCSRFKeyFlag, serverCmd.PersistentFlags().Lookup(serverCSRFKeyFlag),
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
