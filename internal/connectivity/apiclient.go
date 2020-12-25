package connectivity

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Client to call the Plex APIs
type Client struct {
	httpClient *http.Client
	token      string
}

// NewAPIWithHTTPClient creates a new API Client
func NewAPIWithHTTPClient(httpClient *http.Client, token string) *Client {
	return &Client{httpClient: httpClient, token: token}
}

// Call calls ipinfo.io
// Business processing is done in the calling Probe function
func (apiClient *Client) Call() ([]byte, error) {

	req, _ := http.NewRequest("GET", "https://ipinfo.io/", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("token", apiClient.token)
	req.URL.RawQuery = q.Encode()

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
