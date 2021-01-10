package transmission

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Client to call the Transmission APIs
type Client struct {
	Client    *http.Client
	URL       string
	SessionID string
}

// call the specified Transmission API endpoint
// Business processing is done in the calling Probe function
func (client *Client) call(method string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	if client.SessionID, err = client.getSessionID(); err == nil {
		req, _ := http.NewRequest("POST", client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

		if resp, err = client.Client.Do(req); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				body, err = ioutil.ReadAll(resp.Body)
			} else {
				if resp.StatusCode == 409 {
					client.SessionID = resp.Header.Get("X-Transmission-Session-Id")
				}
				err = errors.New(resp.Status)
			}
		}
	}

	return body, err
}

func (client *Client) getSessionID() (string, error) {
	var (
		err       error
		resp      *http.Response
		sessionID = client.SessionID
	)

	if sessionID == "" {
		req, _ := http.NewRequest("POST", client.URL, bytes.NewBufferString("{ \"method\": \"session-get\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

		if resp, err = client.Client.Do(req); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 409 || resp.StatusCode == 200 {
				sessionID = resp.Header.Get("X-Transmission-Session-Id")
			}
		}
		log.WithFields(log.Fields{
			"err":       err,
			"sessionID": sessionID,
		}).Debug("transmission getSessionID")
	}

	return sessionID, err
}
