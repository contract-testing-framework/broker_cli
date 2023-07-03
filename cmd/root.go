package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "signet",
	Short: "The command line interface for the Signet contract testing framework",
	Long:  `The command line interface for the Signet contract testing framework`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
