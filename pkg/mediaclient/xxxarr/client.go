package xxxarr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/go-metrics/client"
	"net/http"
)

// APICaller calls Sonarr / Radarr endpoints
type APICaller interface {
	Get(ctx context.Context, endpoint string, response interface{}) (err error)
	GetURL() string
}

// APIClient is a basic implementation of the APICaller interface
type APIClient struct {
	client.Caller
	URL    string
	APIKey string
}

var _ APICaller = &APIClient{}

// GetURL returns the base URL of the Sonarr / Radarr instance
func (c APIClient) GetURL() string {
	return c.URL
}

// Get calls the specified API endpoint via HTTP GET
func (c APIClient) Get(ctx context.Context, endpoint string, response interface{}) (err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, c.URL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", c.APIKey)

	var resp *http.Response
	resp, err = c.Caller.Do(req)

	if err != nil {
		return fmt.Errorf("call failed: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("call failed: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
