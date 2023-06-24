package cmd

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"net/http"
	"bytes"
	"fmt"
	"encoding/json"

	"github.com/spf13/cobra"
)

var check = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var Type string;
var Branch string;
var ProviderName string;
var Version string;
var ContractFormat string;
var Contract []byte;

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a contract to the broker",
	Long: `Publish a pact contract to the broker.

args:

publish [path to contract] [broker url]


flags:

-t —type         	enum('consumer', 'provider')

-v —version      	service version

-b —branch       	git branch name

-n —provider-name (only for —type 'provider') name of provider service

-c —content-type 	(only for —type 'provider') OAS file type (json or yaml)

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Two arguments are required")
		}

		path := args[0]

		if Type != "consumer" && Type != "provider" {
			if len(Type) == 0 {
				Type = "not set"
			}
			msg := fmt.Sprintf("--type required to be \"consumer\" or \"provider\", --type was %v", Type)
			return errors.New(msg)
		}

		var name string
		if Type == "provider" && len(ProviderName) == 0 {
			return errors.New("Must set --provider-name if --type is \"provider\"")
		}

		if Type == "provider" {
			name = ProviderName
		} else {
			// type is consumer, get consumer name from contract

			type consumer struct{
				Name string `json:"name"`
			}

			type pact struct{
				Consumer consumer `json:"consumer"`
				Interactions interface{} `json:"interactions"`
				MetaData interface{} `json:"metadata"`
				Provider interface{} `json:"provider"`
			}

			contractBytes, err := os.ReadFile(path)
			check(err)

			var contract pact
			err = json.Unmarshal(contractBytes, &contract)
			check(err)

			name = contract.Consumer.Name
		}

		if len(Version) == 0 {
			return errors.New("Must set --version")
		}

		type Body struct{
			ContractType string `json:"contractType"`
			Contract interface{} `json:"contract"`
			ParticipantName string `json:"participantName"`
			ParticipantVersion string `json:"participantVersion"`
			ParticipantBranch string `json:"participantBranch"`
			ContractFormat string `json:"contractFormat"`
		}

		var contractBytes []byte
		var format string
		var contract interface{}

		if path[(len(path) - 4):] == "json" {
			format = "json"

			var err error
			contractBytes, err = os.ReadFile(path)
			check(err)

			err = json.Unmarshal(contractBytes, &contract)
			check(err)

		} else if path[(len(path) - 4):] == "yaml" || path[(len(path) - 3):] == "yml" {
			format = "yaml"

			var err error
			contractBytes, err = os.ReadFile(path)
			check(err)
			contract = string(contractBytes)
			check(err)
		} else {
			return errors.New("Contract must be either JSON or YAML")
		}
		
		brokerURL := args[1]

		requestBody := Body{
			ContractType: Type,
			Contract: contract,
			ParticipantName: name,
			ParticipantVersion: Version,
			ParticipantBranch: Branch,
			ContractFormat: format,
		}

		jsonData, err := json.Marshal(requestBody)
		check(err)

		bodyReader := bytes.NewBuffer(jsonData) // io.Reader interface type

		resp, err := http.Post(brokerURL, "application/json", bodyReader)
		check(err)
		defer resp.Body.Close();

		if resp.StatusCode != 201 {
			type respError struct{
				Error string `json:error`
			}
	
			var respBody respError
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			check(err)

			if respBody.Error == "Participant version already exists" {
				respBody.Error = respBody.Error + "\n\nA new participant version must be set whenever a contract is published."
			}

			fmt.Printf("Status code: %v\n", resp.Status)
			log.Fatal(respBody.Error)
		}

		return nil
	},
}

func init() {
	cmd := exec.Command("git", "rev-parse", "--short=10", "HEAD")
	gitSHA, err := cmd.Output()
	check(err)

	// trim off trailing newline
	gitSHA = gitSHA[:len(gitSHA) - 1]

	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&Type, "type", "t", "", "Type of contract (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&Branch, "branch", "b", "", "Version control branch (optional)")
	publishCmd.Flags().StringVarP(&ProviderName, "provider-name", "n", "", "The name of the provider service (required if --type is \"provider\")")
	publishCmd.Flags().StringVarP(&Version, "version", "v", string(gitSHA), "The version of the service (Defaults to git SHA)")

/*
-v —version (optional)
	- if no version is passed in, use the git branch short SHA
	- API will return 4xx if the consumer version already exists, in that case, log a helpful error msg. (ex. try committing your changes to generate a new git SHA)
*/

}
