package cmd

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func PublishToBroker(brokerURL string, jsonData []byte) error {
	resp, err := http.Post(brokerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		var respBody HttpError
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		}

		if respBody.Error == "Participant version already exists" {
			respBody.Error = respBody.Error + "\n\nA new consumer version must be set whenever a contract is published."
		}

		fmt.Printf("Status code: %v\n", resp.Status)
		log.Fatal(respBody.Error)
	}
	return nil
}

func RegisterEnvWithBroker(brokerURL string, jsonData []byte) error {
	resp, err := http.Post(brokerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		var respBody HttpError
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		}

		fmt.Printf("Status code: %v\n", resp.Status)
		log.Fatal(respBody.Error)
	}
	return nil
}