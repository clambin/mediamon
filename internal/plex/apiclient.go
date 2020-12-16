package plex

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mediamon/internal/version"
	"net/http"
)

type APIClient struct {
	httpClient *http.Client
	url        string
	username   string
	password   string
	authToken  string
}

func NewAPIClient(httpClient *http.Client, url, username, password string) *APIClient {
	return &APIClient{httpClient: httpClient, url: url, username: username, password: password}
}

func (apiClient *APIClient) Call(endpoint string) ([]byte, error) {
	if apiClient.authToken == "" {
		if !apiClient.authenticate() {
			return nil, errors.New("unable to sign in to plex.tv")
		}
	}

	req, _ := http.NewRequest("GET", apiClient.url+endpoint, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Plex-Token", apiClient.authToken)

	resp, err := apiClient.httpClient.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			return ioutil.ReadAll(resp.Body)
		}
		err = errors.New(resp.Status)
	}
	return nil, err
}

func (apiClient *APIClient) authenticate() bool {
	// TODO: there's three different places in the response where the authToken appears.
	// Which is the officially supported version?
	authResponse := struct {
		XMLName             xml.Name `xml:"user"`
		AuthenticationToken string   `xml:"authenticationToken,attr"`
	}{}

	authBody := fmt.Sprintf("user%%5Blogin%%5D=%s&user%%5Bpassword%%5D=%s", apiClient.username, apiClient.password)
	req, _ := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", bytes.NewBufferString(authBody))
	// req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Plex-Product", "mediamon")
	req.Header.Add("X-Plex-Version", version.BuildVersion)
	// TODO: generate UUID?
	req.Header.Add("X-Plex-Client-Identifier", "mediamon-v"+version.BuildVersion)

	resp, err := apiClient.httpClient.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 201 {
			body, _ := ioutil.ReadAll(resp.Body)
			err = xml.Unmarshal(body, &authResponse)
			if err == nil {
				apiClient.authToken = authResponse.AuthenticationToken
				return true
			}
		} else {
			err = errors.New(resp.Status)
		}
	}
	log.Errorf("could not authenticate plex user: %s", err)
	return false
}