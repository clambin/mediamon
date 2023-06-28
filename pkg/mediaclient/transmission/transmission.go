package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Client calls the Transmission APIs
type Client struct {
	HTTPClient *http.Client
	URL        string
}

func NewClient(url string, roundtripper http.RoundTripper) *Client {
	if roundtripper == nil {
		roundtripper = http.DefaultTransport
	}
	return &Client{
		HTTPClient: &http.Client{Transport: &authenticator{next: roundtripper}},
		URL:        url,
	}
}

// GetSessionParameters calls Transmission's "session-get" method. It returns the Transmission instance's configuration parameters
func (client *Client) GetSessionParameters(ctx context.Context) (params SessionParameters, err error) {
	err = client.post(ctx, "session-get", &params)
	if err == nil && params.Result != "success" {
		err = fmt.Errorf("session-get failed: %s", params.Result)
	}
	return
}

// GetSessionStatistics calls Transmission's "session-stats" method. It returns the Transmission instance's session statistics
func (client *Client) GetSessionStatistics(ctx context.Context) (stats SessionStats, err error) {
	err = client.post(ctx, "session-stats", &stats)
	if err == nil && stats.Result != "success" {
		err = fmt.Errorf("session-stats failed: %s", stats.Result)
	}
	return
}

func (client *Client) post(ctx context.Context, method string, response interface{}) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		err = fmt.Errorf("decode: %w", err)
	}
	return err
}
