package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type authDataRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

type authData struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       uint   `json:"expires_in"`
	Interval        uint   `json:"interval"`
}

func loginPrompt(authData *authData) {
	fmt.Printf(
		"You need to log in.\nGo to: %s and enter the code %s\n",
		authData.VerificationURI,
		authData.UserCode,
	)
}

func getAuthData(out chan<- authData) error {

	client := createHTTPClient()
	headers := createUnauthorisedHeaders()
	uri := githubLoginBaseURI + "/device/code"
	relReq := authDataRequest{
		ClientID: "a94fb79ee75e0a953d10",
		Scope:    "repo",
	}

	js, err := json.Marshal(relReq)
	if err != nil {
		return err
	}

	resp, err := client.Post(
		uri,
		bytes.NewReader(js),
		headers,
	)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Couldn't get a device authorisation code")
	}

	defer resp.Body.Close()
	target := authData{}
	err = json.NewDecoder(resp.Body).Decode(&target)
	out <- target
	return nil
}
