package cmd

import (
	"errors"
	"encoding/json"
	// "fmt"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var brokerBaseURL string
var name string

var registerEnvCmd = &cobra.Command{
	Use:   "register-env",
	Short: "register a new deployment environment",
	Long: `register a new deployment environment so that the Signet broker can be informed of which service versions are deployed to which environments.
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(brokerBaseURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		if len(name) == 0 {
			return errors.New("No --name was provided. A value for this flag is required.")
		}

		requestBody := EnvBody{name}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		err = RegisterEnvWithBroker(brokerBaseURL+"/api/environments", jsonData)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(registerEnvCmd)

	registerEnvCmd.Flags().StringVarP(&brokerBaseURL, "broker-url", "u", "", "Scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)")
	registerEnvCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the deployment environment being registered")
}