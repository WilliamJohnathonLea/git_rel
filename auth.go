package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

type authTokenRequest struct {
	ClientID   string `json:"client_id"`
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
}

type authToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (t authToken) AsHTTPHeaderValue() string {
	return fmt.Sprintf("%s %s", t.TokenType, t.AccessToken)
}

func loginPrompt(authData *authData) {
	fmt.Printf(
		"You need to log in.\nGo to: %s and enter the code %s\n",
		authData.VerificationURI,
		authData.UserCode,
	)
	fmt.Println("Press Enter when you have submitted the code")
	fmt.Scanln()
}

func getAuthData(out chan<- authData) error {
	defer close(out)

	client := createHTTPClient()
	headers := createUnauthorisedHeaders()
	uri := githubLoginBaseURI + "/device/code"
	relReq := authDataRequest{
		ClientID: githubClientID,
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

func pollForAuthToken(authData *authData, out chan<- authToken) error {
	defer close(out)
	client := createHTTPClient()
	headers := createUnauthorisedHeaders()
	uri := githubLoginBaseURI + "/oauth/access_token"
	tokenReq := authTokenRequest{
		ClientID:   githubClientID,
		DeviceCode: authData.DeviceCode,
		GrantType:  githubDeviceGrantType,
	}

	js, err := json.Marshal(tokenReq)
	if err != nil {
		return err
	}

	fmt.Println("Authorising...")
	resp, err := client.Post(
		uri,
		bytes.NewReader(js),
		headers,
	)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Couldn't authorise this device")
	}

	defer resp.Body.Close()
	target := authToken{}
	err = json.NewDecoder(resp.Body).Decode(&target)
	if len(target.AccessToken) == 0 {
		return fmt.Errorf("No access token received")
	}

	out <- target

	return nil
}

func saveAuthToken(token *authToken) error {
	js, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("HOME") + "/.git_rel_token.json", js, 0600)
	if err != nil {
		return err
	}

	return nil
}
