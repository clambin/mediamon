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

// GetIdentity calls Plex' /identity endpoint. Mainly useful to get the server's version.
func (c *Client) GetIdentity(ctx context.Context) (identity Identity, err error) {
	err = c.call(ctx, "/identity", &identity)
	return
}

// GetSessions retrieves session information from the server.
func (c *Client) GetSessions(ctx context.Context) (sessions Sessions, err error) {
	err = c.call(ctx, "/status/sessions", &sessions)
	return
}

// GetAuthToken logs into plex.tv and returns the current authToken.
func (c *Client) GetAuthToken(ctx context.Context) (string, error) {
	err := c.authenticate(ctx)
	if err != nil {
		return "", err
	}
	return c.AuthToken, nil
}

func (c *Client) GetLibraries(ctx context.Context) (libraries Libraries, err error) {
	err = c.call(ctx, "/library/sections", &libraries)
	return
}

func (c *Client) GetMovieLibrary(ctx context.Context, key string) (library MovieLibrary, err error) {
	err = c.call(ctx, fmt.Sprintf("/library/sections/%s/all", key), &library)
	return
}

func (c *Client) GetShowLibrary(ctx context.Context, key string) (library ShowLibrary, err error) {
	err = c.call(ctx, fmt.Sprintf("/library/sections/%s/all", key), &library)
	return
}

// call the specified Plex API endpoint
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
