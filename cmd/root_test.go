package cmd

import (
	"testing"
	"bytes"
)

func TestCLIBaseCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs([]string{})
	RootCmd.Execute()

	expected := "The command line interface for the Signet contract testing framework"
	actualOutput := actual.String()

	if actualOutput[:len(expected)] != expected {
		t.Error()
	}
}