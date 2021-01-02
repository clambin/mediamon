package transmission

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// Client to call the Transmission APIs
type Client struct {
	client    *http.Client
	url       string
	sessionID string
}

// call the specified Transmission API endpoint
// Business processing is done in the calling Probe function
func (client *Client) call(method string) ([]byte, error) {
	if client.sessionID == "" {
		if sessionID, err := client.getSessionID(); err == nil {
			client.sessionID = sessionID
		} else {
			return nil, err
		}
	}

	req, _ := http.NewRequest("POST", client.url, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Transmission-Session-Id", client.sessionID)

	resp, err := client.client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if sessionID := resp.Header.Get("X-Transmission-Session-Id"); sessionID != "" {
			client.sessionID = sessionID
		}

		if resp.StatusCode == 200 {
			return ioutil.ReadAll(resp.Body)
		}
		err = errors.New(resp.Status)
	}
	return nil, err
}

func (client *Client) getSessionID() (string, error) {
	req, _ := http.NewRequest("POST", client.url, bytes.NewBufferString("{ \"method\": \"session-get\" }"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Transmission-Session-Id", client.sessionID)

	resp, err := client.client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 409 || resp.StatusCode == 200 {
			return resp.Header.Get("X-Transmission-Session-Id"), nil
		}
	}

	return "", err
}
