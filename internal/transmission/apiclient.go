package transmission

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct {
	client    *http.Client
	url       string
	sessionID string
}

func NewAPIWithHTTPClient(client *http.Client, url string) *Client {
	return &Client{client: client, url: url}
}

func (client *Client) Call(method string) ([]byte, error) {
	if client.sessionID == "" {
		client.sessionID = client.getSessionID()
		if client.sessionID == "" {
			return nil, errors.New("unable to get session ID from Transmission")
		}
	}

	req, _ := http.NewRequest("POST", client.url, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Transmission-Session-Id", client.sessionID)

	resp, err := client.client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			return ioutil.ReadAll(resp.Body)
		}
		err = errors.New(resp.Status)
	}
	return nil, err
}

func (client *Client) getSessionID() string {
	req, _ := http.NewRequest("POST", client.url, bytes.NewBufferString("{ \"method\": \"session-get\" }"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Transmission-Session-Id", client.sessionID)

	resp, err := client.client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 409 || resp.StatusCode == 200 {
			return resp.Header.Get("X-Transmission-Session-Id")
		}
	}

	return ""
}
