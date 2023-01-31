package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// API interface
//
//go:generate mockery --name API
type API interface {
	GetSessionParameters(ctx context.Context) (SessionParameters, error)
	GetSessionStatistics(ctx context.Context) (stats SessionStats, err error)
}

// Client calls the Transmission APIs
type Client struct {
	HTTPClient *http.Client
	URL        string
	SessionID  string
}

var _ API = &Client{}

// GetSessionParameters calls Transmission's "session-get" method. It returns the Transmission instance's configuration parameters
func (client *Client) GetSessionParameters(ctx context.Context) (params SessionParameters, err error) {
	err = client.call(ctx, "session-get", &params)
	if err == nil && params.Result != "success" {
		err = fmt.Errorf("session-get failed: %s", params.Result)
	}
	return
}

// GetSessionStatistics calls Transmission's "session-stats" method. It returns the Transmission instance's session statistics
func (client *Client) GetSessionStatistics(ctx context.Context) (stats SessionStats, err error) {
	err = client.call(ctx, "session-stats", &stats)
	if err == nil && stats.Result != "success" {
		err = fmt.Errorf("session-stats failed: %s", stats.Result)
	}
	return
}

// call the specified Transmission API endpoint
func (client *Client) call(ctx context.Context, method string, response interface{}) (err error) {
	for err == nil {
		var success bool
		success, err = client.post(ctx, method, response)
		if success && err == nil {
			return nil
		}
	}
	return
}

func (client *Client) post(ctx context.Context, method string, response interface{}) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Transmission-Session-Id", client.SessionID)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read: %w", err)
	}

	var success bool
	switch resp.StatusCode {
	case http.StatusOK:
		success = true
		if err = json.Unmarshal(body, response); err != nil {
			err = fmt.Errorf("unmarshal: %w", err)
		}
	case http.StatusConflict:
		client.SessionID = resp.Header.Get("X-Transmission-Session-Id")
	default:
		err = fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	return success, err
}
