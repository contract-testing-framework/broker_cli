package cmd

import (
	"bytes"
	"testing"
	"errors"
	"io/fs"

	utils "github.com/contract-testing-framework/broker_cli/utils"
)

/* ------------- helpers ------------- */

// setup and execute publish command
func callSignetTest(argsAndFlags []string) actualOut {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs(append([]string{"test"}, argsAndFlags...))
	RootCmd.Execute()
	return actualOut{actual.String()}
}

/* ------------- tests ------------- */
func TestSignetTestNoBrokerURL(t *testing.T) {
	flags := []string{
		"--version=version1",
		"--provider-url", "http://localhost:3002",
		"--name", "user_service",
	}
	actual := callSignetTest(flags)
	expected := "Error: No --broker-url was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestSignetTestNoName(t *testing.T) {
	flags := []string{
		"--version=version1",
		"--provider-url", "http://localhost:3002",
		"--broker-url=http://localhost:3000",
	}
	actual := callSignetTest(flags)
	expected := "Error: No --name was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestSignetTestNoProviderURL(t *testing.T) {
	flags := []string{
		"--version=version1",
		"--name", "user_service",
		"--broker-url=http://localhost:3000",
	}
	actual := callSignetTest(flags)
	expected := "Error: No --provider-url was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestSignetCanGetLatestSpec(t *testing.T) {
	realGetNpmPkgRoot := getNpmPkgRoot
	realosWriteFile := osWriteFile
  defer func() { 
		getNpmPkgRoot = realGetNpmPkgRoot
		osWriteFile = realosWriteFile
	}()

	getNpmPkgRoot = func() (string, error) { return "/testDir", nil }

	var specPath string
	var spec []byte
	var rwPermissions fs.FileMode
	osWriteFile = func(name string, data []byte, perm fs.FileMode) error {
		specPath, spec, rwPermissions = name, data, perm
		return errors.New("stop this test here")
	}

	server, req := mockServerForGetSpecsReq200OK(t)
	defer server.Close()

	flags := []string{
		"--version=version1",
		"--name", "user_service",
		"--broker-url", server.URL,
		"--provider-url", "http://localhost:3002",
	}
	actual := callSignetTest(flags)

	t.Run("request has provider query param", func(t *testing.T) {
		if req.URL.Query().Get("provider") != "user_service" {
			t.Error()
		}
	})

	t.Run("calls os.WriteFile with correct arguments", func(t *testing.T) {
		if specPath != "/testDir/specs/spec.json" {
			t.Error()
		}

		if len(spec) == 0 {
			t.Error()
		}

		if rwPermissions != 0666 {
			t.Error()
		}
	})

	t.Run("test stopped at the correct place", func(t *testing.T) {
		expected := "Error: Failed to write specs/spec file: stop this test here"
		actual.startsWith(expected, t)
	})
}

func TestPublishProviderUtilWithoutVersion(t *testing.T) {
	server, reqBody := mockServerForJSONReq201Created[utils.ProviderBody](t)
	defer server.Close()

	path := "../data_test/api-spec.json"
	brokerURL := server.URL
	name := "user_service"
	version := "auto"
	branch := "developement"

	err := utils.PublishProvider(path, brokerURL, name, version, branch)
	if err != nil {
		t.Error()
	}

	t.Run("has correct providerName", func(t *testing.T) {
		if reqBody.ProviderName != "user_service" {
			t.Error()
		}
	})

	t.Run("has a providerVersion", func(t *testing.T) {
		if len(reqBody.ProviderVersion) == 0 {
			t.Error()
		}
	})

	t.Run("has a providerBranch", func(t *testing.T) {
		if len(reqBody.ProviderBranch) == 0 {
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