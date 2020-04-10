package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/francojposa/golang-auth/oauth2-in-action/db"
	"github.com/francojposa/golang-auth/oauth2-in-action/entities/resources"
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

		client := resources.NewClient("example.com")
		fmt.Printf("created Client in app: %q\n", client)

		createdClient, _ := clientRepo.Create(client)

		fmt.Printf("persisted Client in repo: %q\n", createdClient)

		fetchedClient, _ := clientRepo.Get(createdClient.ID)

		fmt.Printf("retrieved Client from repo: %q\n", fetchedClient)

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
