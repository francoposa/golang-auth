package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
		//Run: func(cmd *cobra.Command, args []string) {
		//	// Do stuff here if root command can be used alone
		//},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	serverCmd.PersistentFlags().String("server.host", "", "")
	err := viper.BindPFlag("server.host", serverCmd.PersistentFlags().Lookup("server.host"))
	if err != nil {
		panic(err)
	}

	serverCmd.PersistentFlags().String("server.port", "", "")
	err = viper.BindPFlag("server.port", serverCmd.PersistentFlags().Lookup("server.port"))
	if err != nil {
		panic(err)
	}
}
