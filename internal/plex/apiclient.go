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

// Client to call the Plex APIs
type Client struct {
	Client    *http.Client
	URL       string
	UserName  string
	Password  string
	authToken string
}

// call the specified Plex API endpoint
// Business processing is done in the calling Probe function
func (apiClient *Client) call(endpoint string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	if apiClient.authToken, err = apiClient.authenticate(); err == nil && apiClient.authToken != "" {

		req, _ := http.NewRequest("GET", apiClient.URL+endpoint, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-Plex-Token", apiClient.authToken)

		if resp, err = apiClient.Client.Do(req); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				body, err = ioutil.ReadAll(resp.Body)
			} else {
				err = errors.New(resp.Status)
			}
		}
	}

	return body, err
}

// authenticate logs in to plex.tv and gets an authentication token
// to be used for calls to the Plex server APIs
func (apiClient *Client) authenticate() (string, error) {
	var (
		err       error
		resp      *http.Response
		authToken = apiClient.authToken
	)

	if authToken == "" {
		authBody := fmt.Sprintf("user%%5Blogin%%5D=%s&user%%5Bpassword%%5D=%s",
			apiClient.UserName,
			apiClient.Password,
		)

		req, _ := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", bytes.NewBufferString(authBody))
		req.Header.Add("X-Plex-Product", "mediamon")
		req.Header.Add("X-Plex-Version", version.BuildVersion)
		req.Header.Add("X-Plex-Client-Identifier", "mediamon-v"+version.BuildVersion)

		if resp, err = apiClient.Client.Do(req); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 201 {
				// TODO: there's three different places in the response where the authToken appears.
				// Which is the officially supported version?
				var authResponse struct {
					XMLName             xml.Name `xml:"user"`
					AuthenticationToken string   `xml:"authenticationToken,attr"`
				}

				body, _ := ioutil.ReadAll(resp.Body)
				if err = xml.Unmarshal(body, &authResponse); err == nil {
					authToken = authResponse.AuthenticationToken
				}
			} else {
				err = errors.New(resp.Status)
			}
		}
		log.WithFields(log.Fields{
			"err":       err,
			"authToken": authToken,
		}).Debug("plex authenticate")
	}

	return authToken, err
}
