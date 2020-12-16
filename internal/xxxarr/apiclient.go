package xxxarr

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct {
	client *http.Client
	url    string
	apiKey string
}

func NewAPIWithHTTPClient(client *http.Client, url string, apiKey string) *Client {
	return &Client{client: client, url: url, apiKey: apiKey}
}

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
