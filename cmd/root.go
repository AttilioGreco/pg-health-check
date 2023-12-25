/*
Copyright © 2023 Attilio Greco
*/
package cmd

import (
	"fmt"
	"os"
	"pgHealtchCheck/server"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgHealtchCheck",
	Short: "Web app for PostgreSQL health checks, distinguishing wirte node and read node",
	Long: `This application serves as a web-based health check tool for distingush PostgreSQL node write, and read, usefull with HAProxy, or Kubernetes Healtch Check.

Performs a http request to the /write endpoint to check if the database is in recovery or not.
200 OK --> Node is available for write
403 Forbidden --> Node is not available for write
503 --> DB not reachable

Performs a http request to the /read endpoint to check if the database is reachable or not.
200 OK --> Node is available for read
503 Service Unavailable --> Node is not available for read
`,
	Run: server.Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.pgHealtchCheck.yaml)")
	rootCmd.PersistentFlags().StringP("host", "H", "", "PostgreSQL host")
	rootCmd.PersistentFlags().StringP("port", "p", "", "PostgreSQL port")
	rootCmd.PersistentFlags().StringP("user", "U", "", "PostgreSQL user")
	rootCmd.PersistentFlags().StringP("password", "w", "", "PostgreSQL password")
	rootCmd.PersistentFlags().StringP("dbname", "d", "", "PostgreSQL dbname")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pgHealtchCheck" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pgHealtchCheck")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("Error reading config file")
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
