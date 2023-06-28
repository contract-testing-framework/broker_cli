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

type Body struct {
	ContractType       string      `json:"contractType"`
	Contract           interface{} `json:"contract"`
	ParticipantName    string      `json:"participantName"`
	ParticipantVersion string      `json:"participantVersion"`
	ParticipantBranch  string      `json:"participantBranch"`
	ContractFormat     string      `json:"contractFormat"`
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
			return errors.New("Two arguments are required")
		}
		path := args[0]
		brokerURL := args[1]

		err := ValidFlags()
		if err != nil {
			return err
		}

		format, contract, err := ReadAndUnmarshalContract(path)
		if err != nil {
			return err
		}

		if Type == "consumer" && format != "json" {
			return errors.New("Consumer contracts must be JSON documents")
		}
		
		var name string
		if Type == "provider" {
			name = ProviderName

			if len(Version) != 0 {
				Version = ""
			}
		} else {
			name, err = ConsumerName(path)
			if err != nil {
				return err
			}
		}

		jsonBody, err := CreateRequestBody(Type, contract, name, Version, Branch, format)
		if err != nil {
			return err
		}

		err = publishContract(brokerURL, jsonBody)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cmd := exec.Command("git", "rev-parse", "--short=10", "HEAD")
	gitSHA, err := cmd.Output()
	if err != nil {
		fmt.Printf("Warning: Because this directory is not a git repository, --version cannot default to git commit SHA. --version must be set in order to publish a consumer contract.\n\n")
	}
	// trim off trailing newline
	if len(gitSHA) != 0 {
		gitSHA = gitSHA[:len(gitSHA)-1]
	}

	RootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&Type, "type", "t", "", "Type of contract (\"consumer\" or \"provider\")")
	publishCmd.Flags().StringVarP(&Branch, "branch", "b", "", "Version control branch (optional)")
	publishCmd.Flags().StringVarP(&ProviderName, "provider-name", "n", "", "The name of the provider service (required if --type is \"provider\")")
	publishCmd.Flags().StringVarP(&Version, "version", "v", string(gitSHA), "The version of the service (Defaults to git SHA)")
}

func ValidFlags() error {
	if Type != "consumer" && Type != "provider" {
		if len(Type) == 0 {
			Type = "not set"
		}
		msg := fmt.Sprintf("--type required to be \"consumer\" or \"provider\", --type was %v", Type)
		return errors.New(msg)
	}

	if Type == "provider" && len(ProviderName) == 0 {
		return errors.New("Must set --provider-name if --type is \"provider\"")
	}

	if Type == "consumer" && len(Version) == 0 {
		return errors.New("Must set --version")
	}

	return nil
}

func ReadAndUnmarshalContract(path string) (string, interface{}, error) {
	var contract interface{}

	format := path[len(path)-4:]
	if format != "json" && format != "yaml" && format != ".yml" {
		return "", nil, errors.New("Contract must be either JSON or YAML")
	}

	if format == ".yml" {
		format = "yaml"
	}

	contractBytes, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}

	if format == "json" {
		err = json.Unmarshal(contractBytes, &contract)
	} else {
		contract = string(contractBytes)
	}

	if err != nil {
		return "", nil, err
	}

	return format, contract, nil
}

func ConsumerName(path string) (string, error) {
	contractBytes, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	var contract pact
	err = json.Unmarshal(contractBytes, &contract)

	if err != nil {
		return "", err
	}

	return contract.Consumer.Name, nil
}

func CreateRequestBody(contractType string, contract interface{}, participantName string, participantVersion string, participantBranch string, contractFormat string) ([]byte, error) {
	requestBody := Body{
		ContractType:       contractType,
		Contract:           contract,
		ParticipantName:    participantName,
		ParticipantVersion: participantVersion,
		ParticipantBranch:  participantBranch,
		ContractFormat:     contractFormat,
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
			respBody.Error = respBody.Error + "\n\nA new participant version must be set whenever a contract is published."
		}

		fmt.Printf("Status code: %v\n", resp.Status)
		log.Fatal(respBody.Error)
	}
	return nil
}
