package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/go-metrics/client"
	"net/http"
)

// API interface
//go:generate mockery --name API
type API interface {
	GetSessionParameters(ctx context.Context) (SessionParameters, error)
	GetSessionStatistics(ctx context.Context) (stats SessionStats, err error)
}

// Client calls the Transmission APIs
type Client struct {
	Caller    client.Caller
	URL       string
	SessionID string
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
	var answer bool
	for !answer && err == nil {

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

		var resp *http.Response
		resp, err = client.Caller.Do(req)

		if err != nil {
			break
		}

		switch resp.StatusCode {
		case http.StatusOK:
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(response)
			answer = true
		case http.StatusConflict:
			// Transmission-Session-Id has expired. Get the new one and retry
			client.SessionID = resp.Header.Get("X-Transmission-Session-Id")
		default:
			err = fmt.Errorf("%s", resp.Status)
		}

		_ = resp.Body.Close()
	}

	return
}
