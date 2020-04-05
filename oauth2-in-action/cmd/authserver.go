package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/francojposa/golang-auth/oauth2-in-action/psql"
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
		fmt.Println("authserver called")

		pgConfig := psql.NewDefaultPostgresConfig("OAuth2InAction", "oauth2_in_action")
		db := psql.MustConnect(pgConfig)

		clientRepo := psql.PGClientRepo{Db: db}
		client, _ := clientRepo.GetClient("idtest")

		fmt.Println(client)

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
