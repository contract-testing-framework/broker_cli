package cmd

import (
	"errors"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	client "github.com/contract-testing-framework/broker_cli/client"
	internal "github.com/contract-testing-framework/broker_cli/internal"
)

var Environment string
var Delete bool

var updateDeploymentCmd = &cobra.Command{
	Use:   "update-deployment",
	Short: "notify the broker of a new deployment",
	Long: `notify the broker that a participant version has been deployed to an environment
	
	flags:

	-n --name 					the name of the service

	-v --version        the version of the service

	-e --environment		the name of the environment that the service is deployed to (ex. production)

	-d --delete         the presence of this flag inidicates that the service is no longer deployed to the environment

	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
	
	-i --ignore-config  ingore .signetrc.yaml file if it exists
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		Name = viper.GetString("update-deployment.name")

		if len(BrokerBaseURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		if len(Name) == 0 {
			return errors.New("No --name was provided. A value for this flag is required.")
		}

		if Version == "" || Version == "auto" {
			var err error
			Version, err = internal.SetVersionToGitSha(Version)
			if err != nil {
				return err
			}
		}

		if len(Environment) == 0 {
			return errors.New("No --environment was provided. A value for this flag is required.")
		}

		requestBody := internal.DeploymentBody{
			EnvironmentName: Environment,
			ParticipantName: Name,
			ParticipantVersion: Version,
			Deployed: !Delete,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		err = client.UpdateDeploymentWithBroker(BrokerBaseURL+"/api/participants", jsonData)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateDeploymentCmd)

	updateDeploymentCmd.Flags().StringVarP(&Name, "name", "n", "", "The name of the service which was deployed")
	updateDeploymentCmd.Flags().StringVarP(&Version, "version", "v", "", "The version of the service which was deployed")
	updateDeploymentCmd.Flags().StringVarP(&Environment, "environment", "e", "", "The environment which the service was deployed to")
	updateDeploymentCmd.Flags().BoolVarP(&Delete, "delete", "d", false, "The service is no longer deployed to the environment")
	updateDeploymentCmd.Flags().Lookup("version").NoOptDefVal = "auto"

	viper.BindPFlag("update-deployment.name", updateDeploymentCmd.Flags().Lookup("name"))
}