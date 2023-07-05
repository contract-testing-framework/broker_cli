package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

func teardown() {
	Type = ""
	Path = ""
	BrokerBaseURL = ""
	Branch = ""
	ProviderName = ""
	Version = ""
	ContractFormat = ""
	Contract = []byte{}
	name = ""
	environment = ""
	delete = false
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
	ConsumerBody | ProviderBody | EnvBody | DeploymentBody
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