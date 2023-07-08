package cmd

import (
	"errors"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	utils "github.com/contract-testing-framework/broker_cli/utils"
)

var serviceType string
var contractFormat string
var contract []byte

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a contract or spec to the broker",
	Long: `Publish a consumer contract or provider spec to the broker.

	flags:

	-p --path           the relative path to the contract or spec
	
	-t -—type           the type of service contract (either 'consumer' or 'provider')
	
	-n -—name           canonical name of the provider service (only for —-type 'provider')
	
	-v -—version        service version (only for --type 'consumer', if flag not passed or passed without value, defaults to the git SHA of HEAD)
	
	-b -—branch         git branch name (optional, only for --type 'consumer', defaults to git branch of HEAD)
	
	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
	
	-i --ignore-config  ingore .signetrc.yaml file if it exists
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// get flag values from config file if not passed in on command line
		path = viper.GetString("publish.path")
		serviceType = viper.GetString("publish.type")
		name = viper.GetString("publish.name")

		if len(path) == 0 {
			return errors.New("No --path to a contract/spec was provided. This is a required flag.")
		}

		if len(brokerURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		err := utils.ValidType(serviceType)
		if err != nil {
			return err
		}

		if serviceType == "consumer" {
			err = utils.PublishConsumer(path, brokerURL, version, branch)
			if err != nil {
				return err
			}
		} else {
			err = utils.PublishProvider(path, brokerURL, name, "", "")
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVarP(&path, "path", "p", "", "Relative path from the root directory to the contract or spec file")
	publishCmd.Flags().StringVarP(&serviceType, "type", "t", "", "Type of the participant (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&branch, "branch", "b", "", "git branch name (optional, only for --type 'consumer', defaults to git branch of HEAD)")
	publishCmd.Flags().StringVarP(&name, "name", "n", "", "canonical name of the provider service (only for —-type 'provider')")
	publishCmd.Flags().StringVarP(&version, "version", "v", "", "service version (only for --type 'consumer', if flag not passed or passed without value, defaults to the git SHA of HEAD)")
	publishCmd.Flags().Lookup("version").NoOptDefVal = "auto"
	publishCmd.Flags().Lookup("branch").NoOptDefVal = "auto"

	viper.BindPFlag("publish.path", publishCmd.Flags().Lookup("path"))
	viper.BindPFlag("publish.type", publishCmd.Flags().Lookup("type"))
	viper.BindPFlag("publish.name", publishCmd.Flags().Lookup("name"))
}
