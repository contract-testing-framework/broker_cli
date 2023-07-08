package cmd

import (
	"os"
	"os/exec"
	"errors"
	"log"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	client "github.com/contract-testing-framework/broker_cli/client"
	utils "github.com/contract-testing-framework/broker_cli/utils"
)

const rwPermissions = 0666

var providerURL string

// abstract pkg fn's to enable mocking during testing
var getNpmPkgRoot = utils.GetNpmPkgRoot
var osWriteFile = os.WriteFile

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test that a provider version correctly implements an OpenAPI spec",
	Long: `test that a provider version correctly implements an OpenAPI spec
	
	flags:

	-n --name 					the name of the provider service

	-v --version        the version of the provider service

	-b --branch         Version control branch (optional)

	-s --provider-url   the URL where the provider service is running

	-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
	
	-i --ignore-config  ingore .signetrc.yaml file if it exists
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name = viper.GetString("test.name")
		providerURL = viper.GetString("test.provider-url")

		err := validateTestFlags(brokerURL, name, version, providerURL)
		if err != nil {
			return err
		}

		spec, err := client.GetLatestSpec(brokerURL, name)
		if err != nil {
			return err
		}

		signetRoot, err := getNpmPkgRoot()
		if err != nil {
			return err
		}
		specPath := signetRoot + "/specs/spec.json"
		dreddPath := signetRoot + "/node_modules/dredd"

		err = osWriteFile(specPath, spec, rwPermissions)
		if err != nil {
			return errors.New("Failed to write specs/spec file: " + err.Error())
		}

		testOutput, err := testProvider(dreddPath, specPath, providerURL)

		if err != nil {
			fmt.Println(colorRed + "FAIL" + colorReset + ": Provider test failed - the provider service does not correctly implement the API spec")
			fmt.Println()
			fmt.Println("Breakdown of interactions:")
			testOutput = utils.SliceOutNodeWarnings(testOutput)
			fmt.Println(testOutput)
		} else {
			fmt.Println(colorGreen + "PASS" + colorReset + ": Provider test passed - the provider service correctly implements the API spec")
			fmt.Println()
			fmt.Println("Informing the Signet broker of successful verification...")

			err = utils.PublishProvider(specPath, brokerURL, name, version, branch)
			if err != nil {
				return err
			}

			fmt.Println("Verification results published to Signet broker")
		}

		return nil
	},
}

func validateTestFlags(brokerURL, name, version, providerURL string) error {
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

	if len(providerURL) == 0 {
		return errors.New("No --provider-url was provided. This is a required flag.")
	}

	return nil
}

func testProvider(dreddPath, specPath, providerURL string) (string, error) {
	testCmd := exec.Command("npx", dreddPath, specPath, providerURL, "--loglevel=error")
	stdoutStderr, err := testCmd.CombinedOutput()
	testOutput := string(stdoutStderr)

	if err != nil && len(testOutput) == 0 {
		log.Fatal("Error: failed to execute dredd")
	}

	return testOutput, err
}

func init() {
	RootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the service which was deployed")
	testCmd.Flags().StringVarP(&version, "version", "v", "", "The version of the service which was deployed")
	testCmd.Flags().StringVarP(&branch, "branch", "b", "", "Version control branch (optional)")
	testCmd.Flags().StringVarP(&providerURL, "provider-url", "s", "", "The URL where the provider service is running")
	testCmd.Flags().Lookup("version").NoOptDefVal = "auto"
	testCmd.Flags().Lookup("branch").NoOptDefVal = "auto"

	viper.BindPFlag("test.name", testCmd.Flags().Lookup("name"))
	viper.BindPFlag("test.provider-url", testCmd.Flags().Lookup("provider-url"))
}
