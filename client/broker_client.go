package client

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

/* ---------- client helpers ---------- */

type HttpError struct {
	Error string `json:"error"`
}

func logHTTPErrorThenExit(resp *http.Response) error {
	var respBody HttpError
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	fmt.Printf("Status code: %v\n", resp.Status)
	log.Fatal(respBody.Error)

	return nil
}

/* ---------- client pkg ---------- */

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
	resp, err := http.Post(brokerURL + "/api/environments", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		err = logHTTPErrorThenExit(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateDeploymentWithBroker(brokerURL string, jsonData []byte) error {
	req, err := http.NewRequest(http.MethodPatch, brokerURL + "/api/participants", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
			return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = logHTTPErrorThenExit(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetLatestSpec(brokerURL, name string) (interface{}, error) {
	specURL := brokerURL + "/api/specs?provider=" + name

	resp, err := http.Get(specURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = logHTTPErrorThenExit(resp)
		if err != nil {
			return nil, err
		}
	}

	// move to internal/types.go after dev
	type specResponseBody struct {
		Spec interface{} `json:"spec"`
	}

	var respBody specResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Spec, nil
}