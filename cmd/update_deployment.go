package cmd

import (
	"errors"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	client "github.com/contract-testing-framework/broker_cli/client"
	utils "github.com/contract-testing-framework/broker_cli/utils"
)

var delete bool

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
		name = viper.GetString("update-deployment.name")
		environment = viper.GetString("update-deployment.environment")

		if len(brokerURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		if len(name) == 0 {
			return errors.New("No --name was provided. A value for this flag is required.")
		}

		if version == "" || version == "auto" {
			var err error
			version, err = utils.SetVersionToGitSha(version)
			if err != nil {
				return err
			}
		}

		if len(environment) == 0 {
			return errors.New("No --environment was provided. A value for this flag is required.")
		}

		requestBody := utils.DeploymentBody{
			EnvironmentName: environment,
			ParticipantName: name,
			ParticipantVersion: version,
			Deployed: !delete,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		err = client.UpdateDeploymentWithBroker(brokerURL, jsonData)
		if err != nil {
			return err
		}

		if delete {
			fmt.Println(colorGreen + "Undeployed" + colorReset + " - Signet broker was notified that service version is no longer deployed to the environment")
		} else {
			fmt.Println(colorGreen + "Deployed" + colorReset + " - Signet broker was notified that service version has been deployed to the environment")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateDeploymentCmd)

	updateDeploymentCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the service which was deployed")
	updateDeploymentCmd.Flags().StringVarP(&version, "version", "v", "", "The version of the service which was deployed")
	updateDeploymentCmd.Flags().StringVarP(&environment, "environment", "e", "", "The environment which the service was deployed to")
	updateDeploymentCmd.Flags().BoolVarP(&delete, "delete", "d", false, "The service is no longer deployed to the environment")
	updateDeploymentCmd.Flags().Lookup("version").NoOptDefVal = "auto"

	viper.BindPFlag("update-deployment.name", updateDeploymentCmd.Flags().Lookup("name"))
	viper.BindPFlag("update-deployment.environment", updateDeploymentCmd.Flags().Lookup("environment"))
}