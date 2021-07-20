package mediaclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// TransmissionAPI interface
type TransmissionAPI interface {
	GetVersion(context.Context) (string, error)
	GetStats(context.Context) (int, int, int, int, error)
}

// TransmissionClient calls the Transmission APIs
type TransmissionClient struct {
	Client    *http.Client
	URL       string
	SessionID string
	Options   TransmissionOpts
}

// TransmissionOpts contains options to alter TransmissionClient behaviour
type TransmissionOpts struct {
	PrometheusSummary *prometheus.SummaryVec
}

// GetVersion determines the version of the Transmission server
func (client *TransmissionClient) GetVersion(ctx context.Context) (version string, err error) {
	var (
		resp  []byte
		stats = struct {
			Arguments struct {
				Version string
			}
		}{}
	)

	if resp, err = client.call(ctx, "session-get"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		if err = decoder.Decode(&stats); err == nil {
			version = stats.Arguments.Version
		}
	}

	log.WithFields(log.Fields{
		"err":     err,
		"version": stats.Arguments.Version,
	}).Debug("transmission GetVersion")

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
func (client *TransmissionClient) GetStats(ctx context.Context) (active int, paused int, download int, upload int, err error) {
	var (
		resp  []byte
		stats = struct {
			Arguments struct {
				ActiveTorrentCount int
				PausedTorrentCount int
				UploadSpeed        int
				DownloadSpeed      int
			}
			Result string
		}{}
	)

	if resp, err = client.call(ctx, "session-stats"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		if err = decoder.Decode(&stats); err == nil {
			active = stats.Arguments.ActiveTorrentCount
			paused = stats.Arguments.PausedTorrentCount
			download = stats.Arguments.DownloadSpeed
			upload = stats.Arguments.UploadSpeed
		}
	}

	log.WithFields(log.Fields{
		"err":   err,
		"stats": stats.Arguments,
	}).Debug("transmission GetStats")

	return
}

// call the specified Transmission API endpoint
func (client *TransmissionClient) call(ctx context.Context, method string) (response []byte, err error) {
	var answer bool
	for answer == false && err == nil {

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

		var timer *prometheus.Timer
		if client.Options.PrometheusSummary != nil {
			timer = prometheus.NewTimer(client.Options.PrometheusSummary.WithLabelValues("transmission", method))
		}

		var resp *http.Response
		resp, err = client.Client.Do(req)

		if timer != nil {
			timer.ObserveDuration()
		}

		if err == nil {
			if resp.StatusCode == http.StatusConflict {
				// Transmission-Session-Id has expired. Get the new one and retry
				client.SessionID = resp.Header.Get("X-Transmission-Session-Id")
			} else if resp.StatusCode == http.StatusOK {
				response, err = ioutil.ReadAll(resp.Body)
				answer = true
			} else {
				err = fmt.Errorf("%s", resp.Status)
			}
			_ = resp.Body.Close()
		}
	}
	return
}
