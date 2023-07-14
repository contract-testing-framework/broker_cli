package cmd

import (
	"errors"
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	client "github.com/contract-testing-framework/broker_cli/client"
	utils "github.com/contract-testing-framework/broker_cli/utils"
)

var deployGuardCmd = &cobra.Command{
	Use:   "deploy-guard",
	Short: "check if it is safe to deploy a service version to an environment",
	Long: `check if it is safe to deploy a service version to an environment without breaking any consumers or being broken by an incompatible provider
	
	flags:

	-n --name 					the name of the service
	
	-v --version        the version of the service (defaults to git SHA of HEAD if no value is provided)
	
	-e --environment		the name of the environment that the service is deployed to (ex. production)
	
	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted
	
	-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name = viper.GetString("deploy-guard.name")

		if len(brokerURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}
	
		if len(name) == 0 {
			return errors.New("No --name was provided. This is a required flag.")
		}
	
		if version == "" || version == "auto" {
			var err error
			version, err = utils.SetVersionToGitSha(version)
			if err != nil {
				return err
			}
		}

		if len(environment) == 0 {
			return errors.New("No --environment was provided. This is a required flag.")
		}

		ok, err := client.CheckDeployGuard(brokerURL, name, version, environment)
		if err != nil {
			return err
		}

		if ok {
			cmd.Println(colorGreen + "Safe To Deploy" + colorReset + " - version " + version + " of " + name + " is compatible with all other services in " + environment + " environment")
		} else {
			fmt.Fprintf(os.Stderr, colorRed + "Unsafe to Deploy" + colorReset + " - version " + version + " of " + name + " is incompatible with one or more services in " + environment + " environment\n")
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(deployGuardCmd)

	deployGuardCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the service which was deployed")
	deployGuardCmd.Flags().StringVarP(&version, "version", "v", "auto", "The version of the service which was deployed")
	deployGuardCmd.Flags().StringVarP(&environment, "environment", "e", "", "The environment which the service was deployed to")
	deployGuardCmd.Flags().Lookup("version").NoOptDefVal = "auto"

	viper.BindPFlag("deploy-guard.name", deployGuardCmd.Flags().Lookup("name"))
}