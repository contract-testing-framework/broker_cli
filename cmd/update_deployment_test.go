package cmd

import (
	"bytes"
	"testing"

	utils "github.com/signet-framework/signet-cli/utils"
)

/* ------------- helpers ------------- */

func callUpdateDeployment(argsAndFlags []string) actualOut {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs(append([]string{"update-deployment"}, argsAndFlags...))
	RootCmd.Execute()
	return actualOut{actual.String()}
}

/* ------------- tests ------------- */

func TestUpdateDeploymentNoBrokerURL(t *testing.T) {
	flags := []string{
		"--name", "user_service",
		"--environment", "production",
		"--version=version1",
	}
	actual := callUpdateDeployment(flags)
	expected := "Error: No --broker-url was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestUpdateDeploymentNoName(t *testing.T) {
	flags := []string{
		"--broker-url=http://localhost:3000",
		"--environment", "production",
		"--version=version1",
	}
	actual := callUpdateDeployment(flags)
	expected := "Error: No --name was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestUpdateDeploymentNoEnvironment(t *testing.T) {
	flags := []string{
		"--broker-url=http://localhost:3000",
		"--name", "user_service",
		"--version=version1",
	}
	actual := callUpdateDeployment(flags)
	expected := "Error: No --environment was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestUpdateDeploymentRequest(t *testing.T) {
	server, reqBody := mockServerForJSONReq200OK[utils.DeploymentBody](t)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--name", "user_service",
		"--version=version1",
		"--environment", "production",
	}
	actual := callUpdateDeployment(flags)

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

	t.Run("has correct participantName", func(t *testing.T) {
		if reqBody.ParticipantName != "user_service" {
			t.Error()
		}
	})

	t.Run("has correct participantVersion", func(t *testing.T) {
		if reqBody.ParticipantVersion != "version1" {
			t.Error()
		}
	})

	t.Run("has correct deployed field", func(t *testing.T) {
		if reqBody.Deployed != true {
			t.Error()
		}
	})
	teardown()
}

func TestUpdateDeploymentRequestNoVersion(t *testing.T) {
	server, reqBody := mockServerForJSONReq200OK[utils.DeploymentBody](t)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--name", "user_service",
		"--environment", "production",
	}
	actual := callUpdateDeployment(flags)

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has a participantVersion", func(t *testing.T) {
		if len(reqBody.ParticipantVersion) == 0 {
			t.Error()
		}
	})
	teardown()
}

func TestUpdateDeploymentRequestWithDelete(t *testing.T) {
	server, reqBody := mockServerForJSONReq200OK[utils.DeploymentBody](t)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--name", "user_service",
		"--environment", "production",
		"--version=version1",
		"--delete",
	}
	actual := callUpdateDeployment(flags)

	t.Run("prints nothing to stdout", func(t *testing.T) {
		if actual.actual != "" {
			t.Error()
		}
	})

	t.Run("has correct deployed field", func(t *testing.T) {
		if reqBody.Deployed != false {
			t.Error()
		}
	})
	teardown()
}
