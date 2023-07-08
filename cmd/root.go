package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IgnoreConfig bool
var brokerURL string
var path string
var name string
var version string
var branch string

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "signet",
	Short: "The command line interface for the Signet contract testing framework",
	Long:  `The command line interface for the Signet contract testing framework`,
}

func Execute() {
	readConfigFile()
	brokerURL = viper.GetString("broker-url")

	err := RootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&IgnoreConfig, "ignore-config", "i", false, "ignore config file if present")
	RootCmd.PersistentFlags().StringVarP(&brokerURL, "broker-url", "u", "", "Scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)")

	viper.BindPFlag("broker-url", RootCmd.PersistentFlags().Lookup("broker-url"))
}

func readConfigFile() {
	if IgnoreConfig == false {
		viper.AddConfigPath(".")
		viper.SetConfigName(".signetrc.yaml")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				panic(err)
			}
		}
	}
}