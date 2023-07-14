package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	client "github.com/signet-framework/signet-cli/client"
	utils "github.com/signet-framework/signet-cli/utils"
)

func teardown() {
	serviceType = ""
	path = ""
	brokerURL = ""
	branch = ""
	version = ""
	contractFormat = ""
	contract = []byte{}
	name = ""
	environment = ""
	delete = false
	providerURL = ""
}

type actualOut struct {
	actual string
}

func (ao actualOut) startsWith(expected string, t *testing.T) {
	if len(ao.actual) == 0 {
		fmt.Println("ACTUAL OUTPUT WAS EMPTY")
		t.Error()
	}

	if ao.actual[:len(expected)] != expected {
		fmt.Println("ACTUAL: ")
		fmt.Println(ao.actual)
		t.Error()
	}
}

type requestBody interface {
	utils.ConsumerBody | utils.ProviderBody | utils.EnvBody | utils.DeploymentBody
}

/*
returns a mock server and a pointer to a struct which
will be populated with the request body when a request is made.
used for any requests with a JSON request body, even when contract
is YAML format
*/
func mockServerForJSONReq201Created[T requestBody](t *testing.T) (*httptest.Server, *T) {
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

func mockServerForJSONReq200OK[T requestBody](t *testing.T) (*httptest.Server, *T) {
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

		w.WriteHeader(http.StatusOK)
	}))

	return server, &reqBody
}

func mockServerForGetSpecsReq200OK(t *testing.T) (*httptest.Server, *http.Request) {
	var req http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = *r

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		specBytes, err := os.ReadFile("../data_test/api-spec.json")
		if err != nil {
			t.Error("Failed to load spec for mock response")
		}
		_, err = w.Write(specBytes)
		if err != nil {
			t.Error("Failed to write spec to mock response body")
		}
	}))

	return server, &req
}

func mockServerForDeployGuardReq200OK(t *testing.T, respBody client.DeployGuardResponse) (*httptest.Server, *http.Request) {
	var req http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = *r

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		jsonData, err := json.Marshal(respBody)
		if err != nil {
			t.Error("Failed to encode mock response body")
		}

		_, err = w.Write(jsonData)
		if err != nil {
			t.Error("Failed to write spec to mock response body")
		}
	}))

	return server, &req
}
