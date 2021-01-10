package connectivity

import (
	"errors"
	log "github.com/sirupsen/logrus"
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
	var (
		err  error
		resp *http.Response
		body []byte
	)

	req, _ := http.NewRequest("GET", "https://ipinfo.io/", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("token", apiClient.Token)
	req.URL.RawQuery = q.Encode()

	if resp, err = apiClient.Client.Do(req); err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
		} else {
			err = errors.New(resp.Status)
		}
	}

	log.WithFields(log.Fields{"err": err, "body": body}).Debug("connectivity apiClient call")

	return body, err
}
