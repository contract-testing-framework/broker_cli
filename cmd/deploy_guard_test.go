package cmd

import (
	"testing"
	"bytes"
	"os"
	"os/exec"
	"io/ioutil"

	// utils "github.com/contract-testing-framework/broker_cli/utils"
	client "github.com/contract-testing-framework/broker_cli/client"
)

/* ------------- helpers ------------- */

// setup and execute deploy-guard command
func callDeployGuard(argsAndFlags []string) actualOut {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs(append([]string{"deploy-guard"}, argsAndFlags...))
	RootCmd.Execute()
	return actualOut{actual.String()}
}

/* ------------- tests ------------- */

func TestDeployGuardNoBrokerURL(t *testing.T) {
	flags := []string{
		"--name", "user_service",
		"--environment", "production",
		"--version=version1",
	}
	actual := callDeployGuard(flags)
	expected := "Error: No --broker-url was provided."
	
	actual.startsWith(expected, t)
	teardown()
}

func TestDeployGuardNoName(t *testing.T) {
	flags := []string{
		"--broker-url=http://localhost:3000",
		"--environment", "production",
		"--version=version1",
	}
	actual := callDeployGuard(flags)
	expected := "Error: No --name was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestDeployGuardNoEnvironment(t *testing.T) {
	flags := []string{
		"--broker-url=http://localhost:3000",
		"--name", "user_service",
		"--version=version1",
	}
	actual := callDeployGuard(flags)
	expected := "Error: No --environment was provided."

	actual.startsWith(expected, t)
	teardown()
}

func TestDeployGuardRequest(t *testing.T) {
	respBody := client.DeployGuardResponse{
		Status: true,
		Errors: []client.DeployGuardError{},
	}

	server, req := mockServerForDeployGuardReq200OK(t, respBody)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--name", "user_service",
		"--version=version1",
		"--environment", "production",
	}
	actual := callDeployGuard(flags)

	t.Run("prints 'Safe To Deploy' to stdout", func(t *testing.T) {
		expected := colorGreen + "Safe To Deploy"
		actual.startsWith(expected, t)
	})

	t.Run("request has providerName query param", func(t *testing.T) {
		if req.URL.Query().Get("providerName") != "user_service" {
			t.Error()
		}
	})

	t.Run("request has participantVersion query param", func(t *testing.T) {
		if req.URL.Query().Get("participantVersion") != "version1" {
			t.Error()
		}
	})

	t.Run("request has environmentName query param", func(t *testing.T) {
		if req.URL.Query().Get("environmentName") != "production" {
			t.Error()
		}
	})
	teardown()
}

func TestDeployGuardRequestNoVersion(t *testing.T) {
	respBody := client.DeployGuardResponse{
		Status: true,
		Errors: []client.DeployGuardError{},
	}

	server, req := mockServerForDeployGuardReq200OK(t, respBody)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--name", "user_service",
		"--environment", "production",
	}
	actual := callDeployGuard(flags)

	t.Run("prints 'Safe To Deploy' to stdout", func(t *testing.T) {
		expected := colorGreen + "Safe To Deploy"
		actual.startsWith(expected, t)
	})

	t.Run("request has providerName query param", func(t *testing.T) {
		if req.URL.Query().Get("providerName") != "user_service" {
			t.Error()
		}
	})

	t.Run("request has non-empty participantVersion query param", func(t *testing.T) {
		if len(req.URL.Query().Get("participantVersion")) == 0 {
			t.Error()
		}
	})

	t.Run("request has environmentName query param", func(t *testing.T) {
		if req.URL.Query().Get("environmentName") != "production" {
			t.Error()
		}
	})
	teardown()
}

// to write the following test, see the following articles:
// https://blog.antoine-augusti.fr/2015/12/testing-an-os-exit-scenario-in-golang/
// https://sr-taj.medium.com/how-to-test-methods-that-kill-your-program-in-golang-e3b879185b8a

func TestDeployGuardRequestWhenUnsafe(t *testing.T) {
	respBody := client.DeployGuardResponse{
		Status: false,
		Errors: []client.DeployGuardError{
			client.DeployGuardError{
				Title: "incompatible consumer",
				Details: "service_1 is incompatible with this service as its provider",
			},
		},
	}

	server, _ := mockServerForDeployGuardReq200OK(t, respBody)
	defer server.Close()

	flags := []string{
		"--broker-url", server.URL,
		"--name", "user_service",
		"--version=version1",
		"--environment", "production",
	}

	if os.Getenv("OKAY_TO_EXIT_1") == "true" {
		_ = callDeployGuard(flags)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDeployGuardRequestWhenUnsafe")
	cmd.Env = append(os.Environ(), "OKAY_TO_EXIT_1=true")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Error(err)
	}

	outBytes, _ := ioutil.ReadAll(stdout)
	actual := actualOut{actual: string(outBytes)}

	t.Run("prints 'Unsafe To Deploy' to stdout", func(t *testing.T) {
		expected := colorRed + "Unsafe to Deploy"
		actual.startsWith(expected, t)
	})
	
	err := cmd.Wait()
	t.Run("exits with exit code 1", func(t *testing.T) {
		e, ok := err.(*exec.ExitError)
		if !ok || e.Success() {
			t.Error()
		}
	})

	teardown()
}