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
	HTTPClient *http.Client
	plexAuth   *Authenticator
	URL        string
}

func New(username, password, product, version, url string) *Client {
	auth := &Authenticator{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Username:   username,
		Password:   password,
		AuthURL:    authURL,
		Product:    product,
		Version:    version,
		Next:       http.DefaultTransport,
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
