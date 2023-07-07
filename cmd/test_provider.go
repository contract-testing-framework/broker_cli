package cmd

import (
	// "errors"
	// "encoding/json"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// client "github.com/contract-testing-framework/broker_cli/client"
	// utils "github.com/contract-testing-framework/broker_cli/utils"
)

var ProviderURL string

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test that a provider version correctly implements an OpenAPI spec",
	Long: `test that a provider version correctly implements an OpenAPI spec
	
	flags:

	-n --name 					the name of the provider service

	-v --version        the version of the provider service

	-s --provider-url   the URL where the provider service is running

	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
	
	-i --ignore-config  ingore .signetrc.yaml file if it exists
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shcmd := exec.Command("npm", "root", "-g")
		outBytes, err := shcmd.CombinedOutput()
		if err != nil {
			fmt.Println("47 erred")
			return err
		}
		
		outBytes = outBytes[:len(outBytes) - 1]
		dreddPath := string(outBytes) + "/test_signet_cli/node_modules/dredd"
		fmt.Println(nodeModPath)

		shcmd2 := exec.Command("npx", dreddPath)
		outBytes, err = shcmd2.CombinedOutput()
		if err != nil {
			fmt.Println("53 erred")
			return err
		}

		fmt.Println(string(outBytes))

		// name = viper.GetString("test.name")
		// ProviderURL = viper.GetString("test.provider-url")

		// if len(brokerURL) == 0 {
		// 	return errors.New("No --broker-url was provided. This is a required flag.")
		// }

		// if len(name) == 0 {
		// 	return errors.New("No --name was provided. This is a required flag.")
		// }

		// if version == "" || version == "auto" {
		// 	var err error
		// 	version, err = utils.SetVersionToGitSha(version)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// if len(ProviderURL) == 0 {
		// 	return errors.New("No --provider-url was provided. This is a required flag.")
		// }

		// spec, err := client.GetLatestSpec(brokerURL, name)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(spec)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the service which was deployed")
	testCmd.Flags().StringVarP(&version, "version", "v", "", "The version of the service which was deployed")
	testCmd.Flags().StringVarP(&ProviderURL, "provider-url", "s", "", "The URL where the provider service is running")
	testCmd.Flags().Lookup("version").NoOptDefVal = "auto"

	viper.BindPFlag("test.name", testCmd.Flags().Lookup("name"))
	viper.BindPFlag("test.provider-url", testCmd.Flags().Lookup("provider-url"))
}