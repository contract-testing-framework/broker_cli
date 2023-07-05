package cmd

import (
	"errors"
	"encoding/json"

	"github.com/spf13/cobra"
)

var version string
var environment string
var delete bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "notify the broker of a new deployment",
	Long: `notify the broker that a participant version has been deployed to an environment
	
	flags:

	-n --name 					the name of the service

	-v --version        the version of the service

	-e --environment		the name of the environment that the service is deployed to (ex. production)

	-d --delete         the presence of this flag inidicates that the service is no longer deployed to the environment

	-i --ignore-config  ingore .signetrc.yaml file if it exists
	
	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(BrokerBaseURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		if len(name) == 0 {
			return errors.New("No --name was provided. A value for this flag is required.")
		}

		if version == "" || version == "auto" {
			SetVersionToGitSha()
		}

		if len(environment) == 0 {
			return errors.New("No --environment was provided. A value for this flag is required.")
		}

		requestBody := DeploymentBody{
			EnvironmentName: environment,
			ParticipantName: name,
			ParticipantVersion: version,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		if delete {
			err = DeleteDeploymentFromBroker(BrokerBaseURL+"/api/environments", jsonData)
		} else {
			err = RegisterDeploymentWithBroker(BrokerBaseURL+"/api/environments", jsonData)
			if err != nil {
				return err
			}
		}


		return nil
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the service which was deployed")
	deployCmd.Flags().StringVarP(&version, "version", "v", "", "The version of the service which was deployed")
	deployCmd.Flags().StringVarP(&environment, "environment", "e", "", "The environment which the service was deployed to")
	deployCmd.Flags().BoolVarP(&delete, "delete", "d", false, "The service is no longer deployed to the environment")
	deployCmd.Flags().Lookup("version").NoOptDefVal = "auto"
}