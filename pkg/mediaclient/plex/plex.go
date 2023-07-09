package plex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client calls the Plex APIs
type Client struct {
	URL        string
	HTTPClient *http.Client
	plexAuth   *authenticator
}

func New(username, password, product, version, url string, roundTripper http.RoundTripper) *Client {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	auth := &authenticator{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		username:   username,
		password:   password,
		authURL:    authURL,
		product:    product,
		version:    version,
		next:       roundTripper,
	}

	return &Client{
		URL:        url,
		HTTPClient: &http.Client{Transport: auth},
		plexAuth:   auth,
	}
}

func (c *Client) call(ctx context.Context, endpoint string, response any) error {
	target := c.URL + endpoint
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	req.Header.Add("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	mediaContainer := struct {
		MediaContainer any `json:"MediaContainer"`
	}{MediaContainer: response}

	if err = json.NewDecoder(resp.Body).Decode(&mediaContainer); err != nil {
		err = fmt.Errorf("decode: %w", err)
	}

	return err
}
