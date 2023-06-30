package cmd

import (
	"errors"
	
	"github.com/spf13/cobra"
)

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
		if len(args) < 2 {
			return errors.New("two arguments are required")
		}
		path := args[0]
		brokerBaseUrl := args[1]

		err := ValidType()
		if err != nil {
			return err
		}

		if Type == "consumer" {
			err = PublishConsumer(path, brokerBaseUrl)
			if err != nil {
				return err
			}
		} else {
			err = PublishProvider(path, brokerBaseUrl)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&Type, "type", "t", "", "Type of the participant (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&Branch, "branch", "b", "", "Version control branch (optional)")
	publishCmd.Flags().StringVarP(&ProviderName, "provider-name", "n", "", "The name of the provider service (required if --type is \"provider\")")
	publishCmd.Flags().StringVarP(&Version, "version", "v", "", "The version of the service (Defaults to git SHA)")
	publishCmd.Flags().Lookup("version").NoOptDefVal = "auto"
	publishCmd.Flags().Lookup("branch").NoOptDefVal = "auto"
}
