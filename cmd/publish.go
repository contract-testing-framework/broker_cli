package cmd

import (
	"errors"
	"log"
	"os"
	"net/http"
	"bytes"

	"github.com/spf13/cobra"
)

var check = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a contract to the broker",
	Long: `Publish a pact contract to the broker.

args:

publish [path to contract] [broker url]


flags:

-t —type         	enum('consumer', 'provider')

-v —version      	service version

-b —branch       	git branch name

-n —provider-name (only for —type 'provider') name of provider service

-c —content-type 	(only for —type 'provider') OAS file type (json or yaml)

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Two arguments are required")
		}

		/*
		publish <path to pact> <brokerURL> sends a .json pact to broker
		*/
		path := args[0]
		brokerURL := args[1]
		
		contract, err := os.ReadFile(path)
		check(err)
		bodyReader := bytes.NewReader(contract)

		resp, err := http.Post(brokerURL, "application/json", bodyReader)
		check(err)
		defer resp.Body.Close();


		return nil
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
