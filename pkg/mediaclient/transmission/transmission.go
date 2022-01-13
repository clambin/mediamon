package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/metrics"
	"net/http"
)

// API interface
type API interface {
	GetVersion(context.Context) (string, error)
	GetStats(context.Context) (int, int, int, int, error)
}

// Client calls the Transmission APIs
type Client struct {
	Client    *http.Client
	URL       string
	SessionID string
	Options   Options
}

// Options contains options to alter Client behaviour
type Options struct {
	PrometheusMetrics metrics.PrometheusMetrics
}

// GetVersion determines the version of the Transmission server
func (client *Client) GetVersion(ctx context.Context) (version string, err error) {
	var stats struct {
		Arguments struct {
			Version string
		}
	}

	if err = client.call(ctx, "session-get", &stats); err == nil {
		version = stats.Arguments.Version
	}

	return
}

// GetStats gets torrent & up/download stats from Transmission.
//
// Returns:
//   - active torrents
//   - paused torrents
//   - total download speed
//   - total upload speed
//   - encountered error
func (client *Client) GetStats(ctx context.Context) (active int, paused int, download int, upload int, err error) {
	var stats struct {
		Arguments struct {
			ActiveTorrentCount int
			PausedTorrentCount int
			UploadSpeed        int
			DownloadSpeed      int
		}
		Result string
	}

	if err = client.call(ctx, "session-stats", &stats); err == nil {
		active = stats.Arguments.ActiveTorrentCount
		paused = stats.Arguments.PausedTorrentCount
		download = stats.Arguments.DownloadSpeed
		upload = stats.Arguments.UploadSpeed
	}

	return
}

// call the specified Transmission API endpoint
func (client *Client) call(ctx context.Context, method string, response interface{}) (err error) {
	defer func() {
		client.Options.PrometheusMetrics.ReportErrors(err, "transmission", method)
	}()

	var answer bool
	for answer == false && err == nil {

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

		timer := client.Options.PrometheusMetrics.MakeLatencyTimer("transmission", method)

		var resp *http.Response
		resp, err = client.Client.Do(req)

		if timer != nil {
			timer.ObserveDuration()
		}

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
