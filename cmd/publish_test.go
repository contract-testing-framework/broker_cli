package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

/* ------------- helpers ------------- */

func teardown() {
	Type = ""
	Branch = ""
	ProviderName = ""
	Version = ""
	ContractFormat = ""
	Contract = []byte{}
}

type actualOut struct {
	actual string
}

func (ao actualOut) startsWith(expected string, t *testing.T) {
	if ao.actual[:len(expected)] != expected {
		fmt.Println("ACTUAL: ")
		fmt.Println(ao.actual)
		t.Error()
	}
}

type requestBody interface {
	ConsumerBody | ProviderBody
}

/*
returns a mock server and a pointer to a struct which
will be populated with the request body when a request is made.
used for any requests with a JSON request body, even when contract
is YAML format
*/
func mockServerForJSONReq[T requestBody](t *testing.T) (*httptest.Server, *T) {
	var reqBody T

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Content-Type"))
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			t.Error("Failed to parse request body")
		}

		w.WriteHeader(http.StatusCreated)
	}))

	return server, &reqBody
}

// setup and execute publish command
func callPublish(argsAndFlags []string) actualOut {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs(append([]string{"publish"}, argsAndFlags...))
	RootCmd.Execute()
	return actualOut{actual.String()}
}

/* ------------- tests ------------- */

func TestPublishNoArgs(t *testing.T) {
	actual := callPublish([]string{})
	expected := "Error: two arguments are required"

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishNoType(t *testing.T) {
	args := []string{"../data_test/cons-prov.json", "http://localhost:3000/api/contracts"}
	actual := callPublish(args)
	expected := "Error: --type required to be \"consumer\" or \"provider\", --type was not set"

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishNoProviderName(t *testing.T) {
	args := []string{"../data_test/cons-prov.json", "http://localhost:3000/api/contracts"}
	flags := []string{"--type", "provider"}
	actual := callPublish(append(args, flags...))
	expected := "Error: must set --provider-name if --type is \"provider\""

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishContractDoesNotExist(t *testing.T) {
	args := []string{"../data_test/non-existant.json", "http://localhost:3000/api/contracts"}
	flags := []string{"--type", "consumer"}
	actual := callPublish(append(args, flags...))
	expected := "Error: open ../data_test/non-existant.json: no such file or directory"

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishConsumerWithoutVersionOrBranch(t *testing.T) {
	server, reqBody := mockServerForJSONReq[ConsumerBody](t)
	defer server.Close()

	args := []string{"../data_test/cons-prov.json", server.URL}
	flags := []string{"--type", "consumer"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has a consumerVersion", func(t *testing.T) {
		if len(reqBody.ConsumerVersion) == 0 {
			t.Error()
		}
	})

	t.Run("has a consumerBranch", func(t *testing.T) {
		if len(reqBody.ConsumerBranch) == 0 {
			t.Error()
		}
	})
	teardown()
}

func TestPublishConsumerWithVersion(t *testing.T) {
	server, reqBody := mockServerForJSONReq[ConsumerBody](t)
	defer server.Close()

	args := []string{"../data_test/cons-prov.json", server.URL}
	flags := []string{"--type", "consumer", "--version=version1", "--branch", "main"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct consumerName", func(t *testing.T) {
		if reqBody.ConsumerName != "service_1" {
			t.Error()
		}
	})

	t.Run("has correct consumerVersion", func(t *testing.T) {
		if reqBody.ConsumerVersion != "version1" {
			t.Error()
		}
	})

	t.Run("has correct consumerBranch", func(t *testing.T) {
		if reqBody.ConsumerBranch != "main" {
			t.Error()
		}
	})

	t.Run("has the value of the contract", func(t *testing.T) {
		if reqBody.Contract.Consumer.Name != "service_1" {
			t.Error()
		}
	})
	teardown()
}

func TestPublishProviderWithoutVersion(t *testing.T) {
	server, reqBody := mockServerForJSONReq[ProviderBody](t)
	defer server.Close()

	args := []string{"../data_test/api-spec.json", server.URL}
	flags := []string{"--type", "provider", "--provider-name", "user_service"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct providerName", func(t *testing.T) {
		if reqBody.ProviderName != "user_service" {
			t.Error()
		}
	})

	t.Run("does not have providerVersion", func(t *testing.T) {
		fmt.Println(reqBody.ProviderVersion)
		if len(reqBody.ProviderVersion) != 0 {
			t.Error()
		}
	})

	t.Run("does not have providerBranch", func(t *testing.T) {
		if len(reqBody.ProviderBranch) != 0 {
			t.Error()
		}
	})

	t.Run("has correct specFormat", func(t *testing.T) {
		if reqBody.SpecFormat != "json" {
			t.Error()
		}
	})

	t.Run("has non-nil spec", func(t *testing.T) {
		if reqBody.Spec == nil {
			t.Error()
		}
	})
	teardown()
}

func TestPublishProviderWithVersionAndBranch(t *testing.T) {
	server, reqBody := mockServerForJSONReq[ProviderBody](t)
	defer server.Close()

	args := []string{"../data_test/api-spec.json", server.URL}
	flags := []string{"--type", "provider", "--provider-name", "user_service", "--version=version1", "--branch=main"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct providerVersion", func(t *testing.T) {
		fmt.Println(reqBody.ProviderVersion)
		if reqBody.ProviderVersion != "version1" {
			t.Error()
		}
	})

	t.Run("has correct providerBranch", func(t *testing.T) {
		if reqBody.ProviderBranch != "main" {
			t.Error()
		}
	})

	teardown()
}

func TestPublishProviderYAMLSpec(t *testing.T) {
	server, reqBody := mockServerForJSONReq[ProviderBody](t)
	defer server.Close()

	args := []string{"../data_test/api-spec.yaml", server.URL}
	flags := []string{"--type", "provider", "--provider-name", "user_service"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct providerName", func(t *testing.T) {
		if reqBody.ProviderName != "user_service" {
			t.Error()
		}
	})

	t.Run("does not have providerVersion", func(t *testing.T) {
		fmt.Println(reqBody.ProviderVersion)
		if len(reqBody.ProviderVersion) != 0 {
			t.Error()
		}
	})

	t.Run("does not have providerBranch", func(t *testing.T) {
		if len(reqBody.ProviderBranch) != 0 {
			t.Error()
		}
	})

	t.Run("has correct specFormat", func(t *testing.T) {
		if reqBody.SpecFormat != "yaml" {
			t.Error()
		}
	})

	t.Run("has non-nil spec", func(t *testing.T) {
		if reqBody.Spec == nil {
			t.Error()
		}
	})

	teardown()
}
