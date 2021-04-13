package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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
			Host:                  viper.GetString(pgHostFlag),
			Port:                  viper.GetInt(pgPortFlag),
			Username:              viper.GetString(pgUsernameFlag),
			Password:              viper.GetString(pgPasswordFlag),
			Database:              viper.GetString(pgDatabaseFlag),
			ApplicationName:       viper.GetString(pgApplicationFlag),
			ConnectTimeoutSeconds: viper.GetInt(pgConnectTimeoutFlag),
			SSLMode:               viper.GetString(pgSSLModeFlag),
		}
		sqlxDB := sqlxTools.MustConnect(pgConfig)

		hasher := crypto.NewDefaultArgon2PassHasher()
		userRepo := db.NewPGUserRepo(sqlxDB, hasher)
		userHandler := server.NewUserHandler(userRepo)

		loginRepo := db.NewPGLoginRepo(sqlxDB, userRepo)
		loginURL, err := url.Parse(viper.GetString(uiAppLoginURLFlag))
		if err != nil {
			panic(err)
		}
		loginHandler := server.NewLoginHandler(
			loginRepo,
			userRepo,
			*loginURL,
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
			//router.Put("/", loginHandler.VerifyLogin)
		})

		router.Route("/api/v1/users", func(router chi.Router) {
			router.Post("/", userHandler.Create)
			router.Get("/{id}", userHandler.Get)
		})

		corsAllowedOrigins := viper.GetStringSlice(serverCORSAllowedOriginsFlag)
		fmt.Println(corsAllowedOrigins)
		corsAllowedMethods := viper.GetStringSlice(serverCORSAllowedMethodsFlag)
		fmt.Println(corsAllowedMethods)
		corsAllowCredentials := viper.GetBool(serverCORSAllowCredentialsFlag)
		corsDebug := viper.GetBool(serverCORSDebugFlag)
		corsRouter := cors.New(cors.Options{
			AllowedOrigins:   corsAllowedOrigins,
			AllowedMethods:   corsAllowedMethods,
			AllowCredentials: corsAllowCredentials,
			Debug:            corsDebug,
		}).Handler(router)

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
const serverCORSAllowedOriginsFlag = "server.cors.allowedOrigins"
const serverCORSAllowedMethodsFlag = "server.cors.allowedMethods"
const serverCORSAllowCredentialsFlag = "server.cors.allowCredentials"
const serverCORSDebugFlag = "server.cors.debug"
const serverCSRFSecureFlag = "server.csrf.secure"
const serverCSRFKeyFlag = "server.csrf.key"
const pgHostFlag = "postgres.host"
const pgPortFlag = "postgres.port"
const pgUsernameFlag = "postgres.username"
const pgPasswordFlag = "postgres.password"
const pgDatabaseFlag = "postgres.database"
const pgApplicationFlag = "postgres.application"
const pgConnectTimeoutFlag = "postgres.connectTimeout"
const pgSSLModeFlag = "postgres.sslMode"
const uiAppLoginURLFlag = "ui_app.urls.login"
const uiAppRegisterURLFlag = "ui_app.urls.register"

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
	// HTTP server CORS
	serverCmd.PersistentFlags().String(serverCORSAllowedOriginsFlag, "", "")
	err = viper.BindPFlag(
		serverCORSAllowedOriginsFlag,
		serverCmd.PersistentFlags().Lookup(serverCORSAllowedOriginsFlag),
	)
	serverCmd.PersistentFlags().String(serverCORSAllowedMethodsFlag, "", "")
	err = viper.BindPFlag(
		serverCORSAllowedMethodsFlag,
		serverCmd.PersistentFlags().Lookup(serverCORSAllowedMethodsFlag),
	)
	serverCmd.PersistentFlags().String(serverCORSAllowCredentialsFlag, "", "")
	err = viper.BindPFlag(
		serverCORSAllowCredentialsFlag,
		serverCmd.PersistentFlags().Lookup(serverCORSAllowCredentialsFlag),
	)
	serverCmd.PersistentFlags().String(serverCORSDebugFlag, "", "")
	err = viper.BindPFlag(
		serverCORSDebugFlag,
		serverCmd.PersistentFlags().Lookup(serverCORSDebugFlag),
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
	serverCmd.PersistentFlags().String(pgHostFlag, "", "")
	err = viper.BindPFlag(
		pgHostFlag, serverCmd.PersistentFlags().Lookup(pgHostFlag),
	)
	serverCmd.PersistentFlags().String(pgPortFlag, "", "")
	err = viper.BindPFlag(
		pgPortFlag, serverCmd.PersistentFlags().Lookup(pgPortFlag),
	)
	serverCmd.PersistentFlags().String(pgUsernameFlag, "", "")
	err = viper.BindPFlag(
		pgUsernameFlag, serverCmd.PersistentFlags().Lookup(pgUsernameFlag),
	)
	serverCmd.PersistentFlags().String(pgPasswordFlag, "", "")
	err = viper.BindPFlag(
		pgPasswordFlag, serverCmd.PersistentFlags().Lookup(pgPasswordFlag),
	)
	serverCmd.PersistentFlags().String(pgDatabaseFlag, "", "")
	err = viper.BindPFlag(
		pgDatabaseFlag, serverCmd.PersistentFlags().Lookup(pgDatabaseFlag),
	)
	serverCmd.PersistentFlags().String(pgApplicationFlag, "", "")
	err = viper.BindPFlag(
		pgApplicationFlag, serverCmd.PersistentFlags().Lookup(pgApplicationFlag),
	)
	serverCmd.PersistentFlags().String(pgConnectTimeoutFlag, "", "")
	err = viper.BindPFlag(
		pgConnectTimeoutFlag, serverCmd.PersistentFlags().Lookup(pgConnectTimeoutFlag),
	)
	serverCmd.PersistentFlags().String(pgSSLModeFlag, "", "")
	err = viper.BindPFlag(
		pgSSLModeFlag, serverCmd.PersistentFlags().Lookup(pgSSLModeFlag),
	)

	// Login & Register UI
	serverCmd.PersistentFlags().String(uiAppLoginURLFlag, "", "")
	err = viper.BindPFlag(
		uiAppLoginURLFlag, serverCmd.PersistentFlags().Lookup(uiAppLoginURLFlag),
	)
	serverCmd.PersistentFlags().String(uiAppRegisterURLFlag, "", "")
	err = viper.BindPFlag(
		uiAppRegisterURLFlag, serverCmd.PersistentFlags().Lookup(uiAppRegisterURLFlag),
	)

	if err != nil {
		panic(err)
	}
}
