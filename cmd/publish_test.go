package cmd

import (
	"bytes"
	"testing"
)

/* ------------- helpers ------------- */

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

func TestPublishNoPath(t *testing.T) {
	flags := []string{}
	actual := callPublish(flags)
	expected := "Error: No --path to a contract/spec was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishNoBrokerURL(t *testing.T) {
	flags := []string{
		"--path=../data_test/cons-prov.json",
	}
	actual := callPublish(flags)
	expected := "Error: No --broker-url was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishNoType(t *testing.T) {
	flags := []string{
		"--path=../data_test/cons-prov.json",
		"--broker-url=http://localhost:3000",
	}
	actual := callPublish(flags)
	expected := "Error: --type required to be \"consumer\" or \"provider\", --type was not set"

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishNoProviderName(t *testing.T) {
	flags := []string{
		"--path=../data_test/cons-prov.json",
		"--broker-url=http://localhost:3000",
		"--type", "provider",
	}
	actual := callPublish(flags)
	expected := "Error: must set --provider-name if --type is \"provider\""

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishContractDoesNotExist(t *testing.T) {
	flags := []string{
		"--path=../data_test/non-existant.json",
		"--broker-url=http://localhost:3000",
		"--type", "consumer",
	}
	actual := callPublish(flags)
	expected := "Error: open ../data_test/non-existant.json: no such file or directory"

	actual.startsWith(expected, t)
	teardown()
}

func TestPublishConsumerWithoutVersionOrBranch(t *testing.T) {
	server, reqBody := mockServerForJSONReq[ConsumerBody](t)
	defer server.Close()

	flags := []string{
		"--path=../data_test/cons-prov.json",
		"--broker-url", server.URL,
		"--type", "consumer",
	}
	actual := callPublish(flags)

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

	flags := []string{
		"--path=../data_test/cons-prov.json",
		"--broker-url", server.URL,
		"--type", "consumer",
		"--version=version1",
		"--branch=main",
	}
	actual := callPublish(flags)

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

	flags := []string{
		"--path=../data_test/api-spec.json",
		"--broker-url", server.URL,
		"--type", "provider",
		"--provider-name", "user_service",
	}
	actual := callPublish(flags)

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

	flags := []string{
		"--path=../data_test/api-spec.json",
		"--broker-url", server.URL,
		"--type", "provider",
		"--provider-name", "user_service",
		"--version=version1",
		"--branch=main",
	}
	actual := callPublish(flags)

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct providerVersion", func(t *testing.T) {
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

	flags := []string{
		"--path=../data_test/api-spec.yaml",
		"--broker-url", server.URL,
		"--type", "provider",
		"--provider-name", "user_service",
	}
	actual := callPublish(flags)

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
