package cmd

import (
	"errors"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var path string
var Type string
var Branch string
var ProviderName string
var Version string
var ContractFormat string
var Contract []byte

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a contract or spec to the broker",
	Long: `Publish a consumer contract or provider spec to the broker.

	flags:

	-i --ignore-config  ingore .signetrc.yaml file if it exists
	
	-p --path           the relative path to the contract or spec
	
	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
	
	-t -—type           the type of service contract (either 'consumer' or 'provider')
	
	-n -—provider-name  canonical name of the provider service (only for —-type 'provider')
	
	-v -—version        service version (required for --type 'consumer')
											-—type=consumer: if flag not passed or passed without value, defaults to the git SHA of HEAD
											-—type=provider: if the flag passed without value, defaults to git SHA
	
	-b -—branch         git branch name (optional, defaults to current git branch)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// get flag values from config file if not passed in on command line
		Path = viper.GetString("publish.path")
		Type = viper.GetString("publish.type")
		ProviderName = viper.GetString("publish.provider-name")

		if len(Path) == 0 {
			return errors.New("No --path to a contract/spec was provided. This is a required flag.")
		}

		if len(BrokerBaseURL) == 0 {
			return errors.New("No --broker-url was provided. This is a required flag.")
		}

		err := ValidType()
		if err != nil {
			return err
		}

		if Type == "consumer" {
			err = PublishConsumer(Path, BrokerBaseURL)
			if err != nil {
				return err
			}
		} else {
			err = PublishProvider(Path, BrokerBaseURL)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVarP(&Path, "path", "p", "", "Relative path from the root directory to the contract or spec file")
	publishCmd.Flags().StringVarP(&Type, "type", "t", "", "Type of the participant (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&Branch, "branch", "b", "", "Version control branch (optional)")
	publishCmd.Flags().StringVarP(&ProviderName, "provider-name", "n", "", "The name of the provider service (required if --type is \"provider\")")
	publishCmd.Flags().StringVarP(&Version, "version", "v", "", "The version of the service (Defaults to git SHA)")
	publishCmd.Flags().Lookup("version").NoOptDefVal = "auto"
	publishCmd.Flags().Lookup("branch").NoOptDefVal = "auto"

	viper.BindPFlag("publish.path", publishCmd.Flags().Lookup("path"))
	viper.BindPFlag("publish.type", publishCmd.Flags().Lookup("type"))
	viper.BindPFlag("publish.provider-name", publishCmd.Flags().Lookup("provider-name"))
}
