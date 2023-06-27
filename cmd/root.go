package cmd

import (
		"fmt"
    "os"
    "github.com/spf13/cobra"
    // "github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "broker_cli",
	Short: "A command line interface for the contract broker",
	Long: `A command line interface for the contract broker`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
