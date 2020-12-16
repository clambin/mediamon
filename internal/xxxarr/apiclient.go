package xxxarr

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Client to call the Sonarr/Radarr APIs
type Client struct {
	client *http.Client
	url    string
	apiKey string
}

// NewAPIWithHTTPClient creates a new API Client
func NewAPIWithHTTPClient(client *http.Client, url string, apiKey string) *Client {
	return &Client{client: client, url: url, apiKey: apiKey}
}

// Call calls a specific Sonarr/Radarr API endpoint
// Business processing is done in the calling Probe function
func (client *Client) Call(endpoint string) ([]byte, error) {
	req, _ := http.NewRequest("GET", client.url+endpoint, nil)
	req.Header.Add("X-Api-Key", client.apiKey)

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
