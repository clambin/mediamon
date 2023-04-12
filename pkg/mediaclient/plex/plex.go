package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client calls the Plex APIs
type Client struct {
	HTTPClient *http.Client
	URL        string
	AuthToken  string
	AuthURL    string
	UserName   string
	Password   string
	Product    string
}

func (c *Client) call(ctx context.Context, endpoint string, response any) error {
	if err := c.authenticate(ctx); err != nil {
		return err
	}

	target := c.URL + endpoint
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Plex-Token", c.AuthToken)

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		//return fmt.Errorf("%s: %s", req.URL.Path, resp.Status)
		return fmt.Errorf("%s", resp.Status)
	}

	mediaContainer := struct {
		MediaContainer any `json:"MediaContainer"`
	}{MediaContainer: response}

	if err = json.Unmarshal(body, &mediaContainer); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return err
}
