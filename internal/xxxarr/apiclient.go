package xxxarr

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Client to call the Sonarr/Radarr APIs
type Client struct {
	Client *http.Client
	URL    string
	APIKey string
}

// call the specified Sonarr/Radarr API endpoint
// Business processing is done in the calling Probe function
func (client *Client) call(endpoint string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	req, _ := http.NewRequest("GET", client.URL+endpoint, nil)
	req.Header.Add("X-Api-Key", client.APIKey)

	if resp, err = client.Client.Do(req); err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
		} else {
			err = errors.New(resp.Status)
		}
	}

	return body, err
}
