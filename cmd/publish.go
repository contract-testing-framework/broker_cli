package cmd

import (
	"errors"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IgnoreConfig bool
var Path string
var BrokerBaseURL string
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

arguments:

	publish [path to contract/spec] [broker url]


flags:

-t -—type         	the type of service contract (either 'consumer' or 'provider')

-n -—provider-name 	canonical name of the provider service (only for —-type 'provider')

-v -—version      	service version (required for --type 'consumer')
										-—type=consumer: if flag not passed or passed without value, defaults to the git SHA of HEAD
										-—type=provider: if the flag passed without value, defaults to git SHA

-b -—branch       	git branch name (optional, defaults to current git branch)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		readConfigFile()

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

	publishCmd.Flags().BoolVarP(&IgnoreConfig, "ignore-config", "i", false, "ignore config file if present")
	publishCmd.Flags().StringVarP(&Path, "path", "p", "", "Relative path from the root directory to the contract or spec file")
	publishCmd.Flags().StringVarP(&BrokerBaseURL, "broker-url", "u", "", "Scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)")
	publishCmd.Flags().StringVarP(&Type, "type", "t", "", "Type of the participant (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&Branch, "branch", "b", "", "Version control branch (optional)")
	publishCmd.Flags().StringVarP(&ProviderName, "provider-name", "n", "", "The name of the provider service (required if --type is \"provider\")")
	publishCmd.Flags().StringVarP(&Version, "version", "v", "", "The version of the service (Defaults to git SHA)")
	publishCmd.Flags().Lookup("version").NoOptDefVal = "auto"
	publishCmd.Flags().Lookup("branch").NoOptDefVal = "auto"

	viper.BindPFlag("path", publishCmd.Flags().Lookup("path"))
	viper.BindPFlag("broker-url", publishCmd.Flags().Lookup("broker-url"))
	viper.BindPFlag("type", publishCmd.Flags().Lookup("type"))
	viper.BindPFlag("provider-name", publishCmd.Flags().Lookup("provider-name"))
}

func readConfigFile() {
	if IgnoreConfig == false {
		viper.AddConfigPath(".")
		viper.SetConfigName(".signetrc.yaml")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				panic(err)
			}
		}
	}

	// get flag values from config file if not passed in on command line
	Path = viper.GetString("path")
	BrokerBaseURL = viper.GetString("broker-url")
	Type = viper.GetString("type")
	ProviderName = viper.GetString("provider-name")
}
