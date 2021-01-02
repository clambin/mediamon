package connectivity

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Client to call the Plex APIs
type Client struct {
	Client *http.Client
	Token  string
}

// Call ipinfo.io
// Business processing is done in the calling Probe function
func (apiClient *Client) call() ([]byte, error) {

	req, _ := http.NewRequest("GET", "https://ipinfo.io/", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("token", apiClient.Token)
	req.URL.RawQuery = q.Encode()

	resp, err := apiClient.Client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			return ioutil.ReadAll(resp.Body)
		}
		err = errors.New(resp.Status)
	}
	return nil, err
}
