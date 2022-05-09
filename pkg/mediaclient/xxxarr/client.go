package xxxarr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type apiClient struct {
	HTTPClient  *http.Client
	URL         string
	APIKey      string
	options     Options
	application string
}

// Get calls the specified API endpoint via HTTP GET
func (c apiClient) Get(ctx context.Context, endpoint string, response interface{}) (err error) {
	defer func() {
		c.options.PrometheusMetrics.ReportErrors(err, c.application, endpoint)
	}()

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, c.URL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", c.APIKey)

	timer := c.options.PrometheusMetrics.MakeLatencyTimer(c.application, endpoint)
	var resp *http.Response
	resp, err = c.HTTPClient.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}

	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return json.NewDecoder(resp.Body).Decode(response)
}
