package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)


/* ------------- helpers ------------- */

type actualOut struct{
	actual string
}

func (ao actualOut) startsWith(expected string, t *testing.T) {
	if ao.actual[:len(expected)] != expected {
		t.Error()
	}
}

/*
	returns a mock server and a pointer to a struct which 
	will be populated with the request body when a request is made.
	used for any requests with a JSON request body, even when contract
	is YAML format
*/
func mockServerForJSONReq(t *testing.T) (*httptest.Server, *Body) {
	var reqBody Body

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
	expected := "Error: Two arguments are required"

	actual.startsWith(expected, t)
}

func TestPublishNoType(t *testing.T) {
	args := []string{"../data_test/cons-prov.json", "http://localhost:3000/api/contracts"}
	actual := callPublish(args)
	expected := "Error: --type required to be \"consumer\" or \"provider\", --type was not set"

	actual.startsWith(expected, t)
}

func TestPublishNoProviderName(t *testing.T) {
	args := []string{"../data_test/cons-prov.json", "http://localhost:3000/api/contracts"}
	flags := []string{"--type", "provider"}
	actual := callPublish(append(args, flags...))
	expected := "Error: Must set --provider-name if --type is \"provider\""

	actual.startsWith(expected, t)
}

func TestPublishContractDoesNotExist(t *testing.T) {
	args := []string{"../data_test/non-existant.json", "http://localhost:3000/api/contracts"}
	flags := []string{"--type", "consumer"}
	actual := callPublish(append(args, flags...))
	expected := "Error: open ../data_test/non-existant.json: no such file or directory"

	actual.startsWith(expected, t)
}

func TestPublishConsumerContract(t *testing.T) {
	server, reqBody := mockServerForJSONReq(t)
	defer server.Close()

	args := []string{"../data_test/cons-prov.json", server.URL}
	flags := []string{"--type", "consumer", "--branch", "main"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct contractType", func(t *testing.T) {
		if reqBody.ContractType != "consumer" {
			t.Error()
		}
	})

	t.Run("has correct participantName", func(t *testing.T) {
		if reqBody.ParticipantName != "service_1" {
			t.Error()
		}
	})

	t.Run("has a participantVersion", func(t *testing.T) {
		if len(reqBody.ParticipantVersion) == 0 {
			t.Error()
		}
	})

	t.Run("has correct participantBranch", func(t *testing.T) {
		if reqBody.ParticipantBranch != "main" {
			t.Error()
		}
	})

	t.Run("has correct contractFormat", func(t *testing.T) {
		if reqBody.ContractFormat != "json" {
			t.Error()
		}
	})

	t.Run("has non-null contract", func(t *testing.T) {
		if reqBody.Contract == nil {
			t.Error()
		}
	})
}

func TestPublishProviderJSONSpec(t *testing.T) {
	server, reqBody := mockServerForJSONReq(t)
	defer server.Close()

	args := []string{"../data_test/api-spec.json", server.URL}
	flags := []string{"--type", "provider", "--provider-name", "user_service", "--branch", "main"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct contractType", func(t *testing.T) {
		if reqBody.ContractType != "provider" {
			t.Error()
		}
	})

	t.Run("has correct participantName", func(t *testing.T) {
		if reqBody.ParticipantName != "user_service" {
			t.Error()
		}
	})

	t.Run("does not have participantVersion", func(t *testing.T) {
		if len(reqBody.ParticipantVersion) != 0 {
			t.Error()
		}
	})

	t.Run("has correct participantBranch", func(t *testing.T) {
		if reqBody.ParticipantBranch != "main" {
			t.Error()
		}
	})

	t.Run("has correct contractFormat", func(t *testing.T) {
		if reqBody.ContractFormat != "json" {
			t.Error()
		}
	})

	t.Run("has non-null contract", func(t *testing.T) {
		if reqBody.Contract == nil {
			t.Error()
		}
	})
}

func TestPublishProviderYAMLSpec(t *testing.T) {
	server, reqBody := mockServerForJSONReq(t)
	defer server.Close()

	args := []string{"../data_test/api-spec.yaml", server.URL}
	flags := []string{"--type", "provider", "--provider-name", "user_service", "--branch", "main"}
	actual := callPublish(append(args, flags...))

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct contractType", func(t *testing.T) {
		if reqBody.ContractType != "provider" {
			t.Error()
		}
	})

	t.Run("has correct participantName", func(t *testing.T) {
		if reqBody.ParticipantName != "user_service" {
			t.Error()
		}
	})

	t.Run("does not have participantVersion", func(t *testing.T) {
		if len(reqBody.ParticipantVersion) != 0 {
			t.Error()
		}
	})

	t.Run("has correct participantBranch", func(t *testing.T) {
		if reqBody.ParticipantBranch != "main" {
			t.Error()
		}
	})

	t.Run("has correct contractFormat", func(t *testing.T) {
		if reqBody.ContractFormat != "yaml" {
			t.Error()
		}
	})

	t.Run("has non-null contract", func(t *testing.T) {
		if reqBody.Contract == nil {
			t.Error()
		}
	})
}
