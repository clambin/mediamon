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
	req, _ := http.NewRequest("GET", client.URL+endpoint, nil)
	req.Header.Add("X-Api-Key", client.APIKey)

	resp, err := client.Client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			return ioutil.ReadAll(resp.Body)
		}
		err = errors.New(resp.Status)
	}
	return nil, err
}
