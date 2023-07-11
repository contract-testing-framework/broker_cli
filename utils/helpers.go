package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	client "github.com/contract-testing-framework/broker_cli/client"
)

func ValidType(serviceType string) error {
	if serviceType != "consumer" && serviceType != "provider" {
		if len(serviceType) == 0 {
			serviceType = "not set"
		}
		msg := fmt.Sprintf("--type required to be \"consumer\" or \"provider\", --type was %v", serviceType)
		return errors.New(msg)
	}
	return nil
}

func LoadContract(path string) (contract Pact, err error) {
	contractBytes, err := os.ReadFile(path)
	if err != nil {
		return Pact{}, err
	}

	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return Pact{}, err
	}
	return
}

func LoadSpec(path string) (spec interface{}, format string, err error) {
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

func CreateConsumerRequestBody(contract Pact, consumerName string, consumerVersion string, consumerBranch string) ([]byte, error) {

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

func CreateProviderRequestBody(spec interface{}, providerName string, providerVersion string, providerBranch string, specFormat string) ([]byte, error) {
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

func SetVersionToGitSha(version string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short=10", "HEAD")
	gitSHA, err := cmd.Output()
	if err != nil {
		return "", errors.New("because this directory is not a git repository, --version cannot default to git commit SHA. --version must be set for this command.")
	}
	if len(gitSHA) != 0 {
		gitSHA = gitSHA[:len(gitSHA)-1]
	}

	return string(gitSHA), nil
}

func SetBranchToCurrentGit(branch string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	currentBranch, err := cmd.Output()
	if err != nil {
		return "", errors.New("because this directory is not a git repository, --branch cannot default to current git branch")
	}
	if len(currentBranch) != 0 {
		currentBranch = currentBranch[:len(currentBranch)-1]
	}

	return string(currentBranch), nil
}

func PublishConsumer(path string, brokerURL string, version, branch string) error {
	if branch == "auto" || (branch == "" && (version == "auto" || version == "")) {
		var err error
		branch, err = SetBranchToCurrentGit(branch)
		if err != nil {
			return err
		}
	}

	if version == "" || version == "auto" {
		var err error
		version, err = SetVersionToGitSha(version)
		if err != nil {
			return err
		}
	}

	contract, err := LoadContract(path)
	if err != nil {
		return err
	}

	consumerName := contract.Consumer.Name

	if len(consumerName) == 0 {
		return errors.New("consumer contract does not have a consumer name")
	}

	requestBody, err := CreateConsumerRequestBody(contract, consumerName, version, branch)
	if err != nil {
		return err
	}

	err = client.PublishToBroker(brokerURL+"/api/contracts", requestBody)
	if err != nil {
		return err
	}

	return nil
}

func PublishProvider(path string, brokerURL string, ProviderName, version, branch string) error {
	if len(ProviderName) == 0 {
		return errors.New("must set --name if --type is \"provider\"")
	}

	if branch == "auto" || (branch == "" && version == "auto") {
		var err error
		branch, err = SetBranchToCurrentGit(branch)
		if err != nil {
			return err
		}
	}

	if version == "auto" {
		var err error
		version, err = SetVersionToGitSha(version)
		if err != nil {
			return err
		}
	}

	spec, specFormat, err := LoadSpec(path)
	if err != nil {
		return err
	}

	requestBody, err := CreateProviderRequestBody(spec, ProviderName, version, branch, specFormat)
	if err != nil {
		return err
	}

	err = client.PublishToBroker(brokerURL+"/api/specs", requestBody)
	if err != nil {
		return err
	}

	return nil
}

func SliceOutNodeWarnings(str string) string {
	re := regexp.MustCompile(`(?s)\(node(.+)warning was created\)\n`)
	return re.ReplaceAllString(str, "")
}

func GetNpmPkgRoot() (string, error) {
	shcmd := exec.Command("npm", "root", "-g")
	stdoutStderr, err := shcmd.CombinedOutput()
	if err != nil {
		return "", errors.New("Could not find npm root")
	}

	if len(stdoutStderr) < 1 {
		return "", errors.New("npm root path was empty string")
	}

	pkgRoot := string(stdoutStderr[:len(stdoutStderr)-1]) + "/signet-cli"

	return pkgRoot, nil
}

func CreatePact(stubsPath string, pactPath string, consumerName string, providerName string) (error, bool) {

	pact := CreateDefaultPact(pactPath, consumerName, providerName)
	matchPaths, err := GetMatchPaths(stubsPath)

	if err != nil {
		return err, false
	}

	interactions, err := createInteractions(matchPaths)
	pact["interactions"] = interactions

	err = WritePact(pact, pactPath)

	if err != nil {
		return err, false
	}

	return err, true
}

func GetMatchPaths(stubsPath string) ([]string, error) {
	matchPaths := []string{}

	filepath.WalkDir(stubsPath, func(path string, d fs.DirEntry, err error) error {
		dir := filepath.Dir(path)
		parent := filepath.Base(dir)

		if parent == "matches" {
			matchPaths = append(matchPaths, path)
		}
		return err
	})
	return matchPaths, nil
}

func WritePact(pact map[string]interface{}, pactPath string) error {
	CreatePactDir(pactPath)
	//fmt.Print(pact)
	file, _ := json.MarshalIndent(pact, "", " ")

	err := os.WriteFile(pactPath, file, 0644)

	return err
}

func CreatePactDir(pactDir string) error {
	err := os.MkdirAll(filepath.Dir(pactDir), os.ModePerm)

	return err
}

func createInteractions(matchPaths []string) ([]map[string]interface{}, error) {
	interactions := []map[string]interface{}{}

	for _, matchPath := range matchPaths {
		matchBytes, err := os.ReadFile(matchPath)
		if err != nil {
			return []map[string]interface{}{}, err
		}

		match := map[string]interface{}{}
		err = json.Unmarshal(matchBytes, &match)

		if err != nil {
			return []map[string]interface{}{}, err
		}

		interaction := map[string]interface{}{}

		request := match["request"].(map[string]any)
		response := match["response"].(map[string]any)

		interaction["description"] = fmt.Sprintf("%s %s %.0f", request["method"], request["path"], response["statusCode"])

		/*
			var requestBody map[string]interface{}

			if request["body"] != "" {
				err = json.Unmarshal([]byte(request["body"].(string)), &requestBody)
			}
		*/

		/*
			var requestQuery string

			if len(request["query"].(map[string]interface{})) > 0 {
				params := url.Values{}

				for key, values := range request["query"].(map[string]interface{}) {
					for _, value := range values.([]interface{}) {
						params.Add(key, value.(string))
					}
				}
				requestQuery = url.Values(params).Encode()
			}*/

		interaction["request"] = map[string]interface{}{
			"method":  request["method"],
			"path":    request["path"],
			"body":    request["body"],
			"query":   request["query"],
			"headers": request["headers"],
		}

		interaction["response"] = map[string]interface{}{
			"status":  response["statusCode"],
			"headers": response["headers"],
			"body":    response["body"],
		}

		interactions = append(interactions, interaction)
	}
	return interactions, nil
}

func CreateDefaultPact(pactPath string, consumerName string, providerName string) (contract map[string]interface{}) {
	return map[string]interface{}{
		"consumer": map[string]interface{}{
			"name": consumerName,
		},
		"provider": map[string]interface{}{
			"name": providerName,
		},
		"interactions": nil,
		"metadata": map[string]interface{}{
			"pactSpecification": map[string]interface{}{
				"version": "3.0.0",
			},
		},
	}
}
