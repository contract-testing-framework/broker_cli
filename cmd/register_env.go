package cmd

import (
	"encoding/json"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	client "github.com/signet-framework/signet-cli/client"
	utils "github.com/signet-framework/signet-cli/utils"
)

var registerEnvCmd = &cobra.Command{
	Use:   "register-env",
	Short: "register a new deployment environment",
	Long: `register a new deployment environment so that the Signet broker can be informed of which service versions are deployed to which environments.
	
	flags:

	-e --environment    the name of the deployment environment being registered (ex. production)

	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted

	-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		environment = viper.GetString("register-env.environment")

		if len(brokerURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		if len(environment) == 0 {
			return errors.New("No --environment was provided. A value for this flag is required.")
		}

		requestBody := utils.EnvBody{EnvironmentName: environment}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		err = client.RegisterEnvWithBroker(brokerURL, jsonData)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(registerEnvCmd)

	registerEnvCmd.Flags().StringVarP(&environment, "environment", "e", "", "The name of the deployment environment being registered")

	viper.BindPFlag("register-env.environment", registerEnvCmd.Flags().Lookup("environment"))
}
