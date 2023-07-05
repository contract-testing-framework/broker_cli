package cmd

import (
	"testing"
	"bytes"
)

/* ------------- helpers ------------- */

// setup and execute register-env command
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
		"--name=production",
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
	expected := "Error: No --name was provided."

	actual.startsWith(expected, t)
	teardown()
}