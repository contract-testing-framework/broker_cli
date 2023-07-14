package cmd

import (
	"bytes"
	"testing"

	utils "github.com/signet-framework/signet-cli/utils"
)

/* ------------- helpers ------------- */

func callRegisterEnv(argsAndFlags []string) actualOut {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs(append([]string{"register-env"}, argsAndFlags...))
	RootCmd.Execute()
	return actualOut{actual.String()}
}

/* ------------- tests ------------- */

func TestRegisterEnvNoBrokerURL(t *testing.T) {
	flags := []string{
		"--environment=production",
	}
	actual := callRegisterEnv(flags)
	expected := "Error: No --broker-url was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestRegisterEnvNoName(t *testing.T) {
	flags := []string{
		"--broker-url=http://localhost:3000",
	}
	actual := callRegisterEnv(flags)
	expected := "Error: No --environment was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestRegisterEnvRequest(t *testing.T) {
	server, reqBody := mockServerForJSONReq201Created[utils.EnvBody](t)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--environment=production",
	}
	actual := callRegisterEnv(flags)

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct environmentName", func(t *testing.T) {
		if reqBody.EnvironmentName != "production" {
			t.Error()
		}
	})
	teardown()
}
