package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type consumer struct {
	Name string `json:"name"`
}

type pact struct {
	Consumer     consumer    `json:"consumer"`
	Interactions interface{} `json:"interactions"`
	MetaData     interface{} `json:"metadata"`
	Provider     interface{} `json:"provider"`
}

type ConsumerBody struct {
	Contract        pact   `json:"contract"`
	ConsumerName    string `json:"consumerName"`
	ConsumerVersion string `json:"consumerVersion"`
	ConsumerBranch  string `json:"consumerBranch"`
}

type ProviderBody struct {
	Spec            interface{} `json:"spec"`
	ProviderName    string      `json:"providerName"`
	ProviderVersion string      `json:"providerVersion"`
	ProviderBranch  string      `json:"providerBranch"`
	SpecFormat      string      `json:"specFormat"`
}

type httpError struct {
	Error string `json:"error"`
}

var Type string
var Branch string
var ProviderName string
var Version string
var ContractFormat string
var Contract []byte

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a contract to the broker",
	Long: `Publish a pact contract to the broker.

args:

publish [path to contract] [broker url]


flags:

-t -—type         	the type of service contract (either 'consumer' or 'provider')

-b -—branch       	git branch name (optional)

-v -—version      	version of service (only for --type 'consumer', defaults to SHA of git commit)

-n -—provider-name 	identifier key for provider service (only for —-type 'provider')
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
			if Version == "" || Version == "auto" {
				setVersionToGitSha()
			}

			contract, err := loadContract(path)
			if err != nil {
				return err
			}

			consumerName := contract.Consumer.Name

			if len(consumerName) == 0 {
				return errors.New("consumer contract does not have a consumer name")
			}

			requestBody, err := createConsumerRequestBody(contract, consumerName, Version, Branch)
			if err != nil {
				return err
			}

			err = publishContract(brokerBaseUrl+"/api/contracts", requestBody)
			if err != nil {
				return err
			}
			// Else it is a provider
		} else {
			if len(ProviderName) == 0 {
				return errors.New("must set --provider-name if --type is \"provider\"")
			}

			if Version == "auto" {
				setVersionToGitSha()
			}

			spec, specFormat, err := loadSpec(path)
			if err != nil {
				return err
			}

			requestBody, err := createProviderRequestBody(spec, ProviderName, Version, Branch, specFormat)
			if err != nil {
				return err
			}
			err = publishContract(brokerBaseUrl+"/api/specs", requestBody)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&Type, "type", "t", "", "Type of contract (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&Branch, "branch", "b", "", "Version control branch (optional)")
	publishCmd.Flags().StringVarP(&ProviderName, "provider-name", "n", "", "The name of the provider service (required if --type is \"provider\")")
	publishCmd.Flags().StringVarP(&Version, "version", "v", "", "The version of the service (Defaults to git SHA)")
	publishCmd.Flags().Lookup("version").NoOptDefVal = "auto"
}

func ValidType() error {
	if Type != "consumer" && Type != "provider" {
		if len(Type) == 0 {
			Type = "not set"
		}
		msg := fmt.Sprintf("--type required to be \"consumer\" or \"provider\", --type was %v", Type)
		return errors.New(msg)
	}
	return nil
}

func loadContract(path string) (contract pact, err error) {
	contractBytes, err := os.ReadFile(path)
	if err != nil {
		return pact{}, err
	}

	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return pact{}, err
	}
	return
}

func loadSpec(path string) (spec interface{}, format string, err error) {
	format = path[len(path)-4:]
	if format != "json" && format != "yaml" && format != ".yml" {
		return nil, "", errors.New("spec must be either JSON or YAML")
	}

	if format == ".yml" {
		format = "yaml"
	}

	specBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	if format == "json" {
		err = json.Unmarshal(specBytes, &spec)
	} else {
		spec = string(specBytes)
	}

	if err != nil {
		return nil, "", err
	}

	return
}

func createConsumerRequestBody(contract pact, consumerName string, consumerVersion string, consumerBranch string) ([]byte, error) {

	requestBody := ConsumerBody{
		Contract:        contract,
		ConsumerName:    consumerName,
		ConsumerVersion: consumerVersion,
		ConsumerBranch:  consumerBranch,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func createProviderRequestBody(spec interface{}, providerName string, providerVersion string, providerBranch string, specFormat string) ([]byte, error) {
	requestBody := ProviderBody{
		Spec:            spec,
		ProviderName:    providerName,
		ProviderVersion: providerVersion,
		ProviderBranch:  providerBranch,
		SpecFormat:      specFormat,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func publishContract(brokerURL string, jsonData []byte) error {
	resp, err := http.Post(brokerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		var respBody httpError
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		}

		if respBody.Error == "Participant version already exists" {
			respBody.Error = respBody.Error + "\n\nA new consumer version must be set whenever a contract is published."
		}

		fmt.Printf("Status code: %v\n", resp.Status)
		log.Fatal(respBody.Error)
	}
	return nil
}

func setVersionToGitSha() error {
	cmd := exec.Command("git", "rev-parse", "--short=10", "HEAD")
	gitSHA, err := cmd.Output()
	if err != nil {
		return errors.New("because this directory is not a git repository, --version cannot default to git commit SHA. --version must be set in order to publish a consumer contract")
	}
	if len(gitSHA) != 0 {
		gitSHA = gitSHA[:len(gitSHA)-1]
	}

	Version = string(gitSHA)
	return nil
}
