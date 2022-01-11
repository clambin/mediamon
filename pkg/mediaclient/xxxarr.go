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

// XXXArrAPI interface
type XXXArrAPI interface {
	GetVersion(ctx context.Context) (string, error)
	GetCalendar(ctx context.Context) (int, error)
	GetQueue(ctx context.Context) (int, error)
	GetMonitored(ctx context.Context) (int, int, error)
	GetApplication() string
	GetURL() string
}

// XXXArrClient calls the Sonarr/Radarr APIs.  Application specifies whether this is a
// Sonarr ("sonarr") or Radarr ("radarr") server.  XXXArrClient will panic if Application
// contains any other values.
type XXXArrClient struct {
	Client      *http.Client
	URL         string
	APIKey      string
	Application string
	Options     XXXArrOpts
}

// XXXArrOpts contains options to alter XXXArrClient behaviour
type XXXArrOpts struct {
	PrometheusSummary *prometheus.SummaryVec
}

// GetApplication returns the client's configured application
func (client *XXXArrClient) GetApplication() string {
	return client.Application
}

// GetURL returns the server URL
func (client *XXXArrClient) GetURL() string {
	return client.URL
}

// GetVersion retrieves the version of the Sonarr/Radarr server
func (client *XXXArrClient) GetVersion(ctx context.Context) (version string, err error) {
	var resp []byte
	resp, err = client.call(ctx, "/api/v3/system/status")

	if err == nil {
		var stats struct {
			Version string
		}

		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)

		if err == nil {
			version = stats.Version
		}
	}

	log.WithError(err).WithFields(log.Fields{
		"application": client.Application,
		"version":     version,
	}).Debug("xxxarr Version")

	return
}

// GetCalendar retrieves the number of upcoming movies/series airing today and tomorrow
func (client *XXXArrClient) GetCalendar(ctx context.Context) (calendar int, err error) {
	var resp []byte
	resp, err = client.call(ctx, "/api/v3/calendar")

	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		var stats []struct{ HasFile bool }
		if err = decoder.Decode(&stats); err == nil {
			calendar = 0
			for _, stat := range stats {
				if stat.HasFile == false {
					calendar++
				}
			}
		}
	}

	log.WithError(err).WithFields(log.Fields{
		"application": client.Application,
		"calendar":    calendar,
	}).Debug("xxxarr getCalendar")

	return
}

// GetQueue retrieves how many movies/series are currently downloading
func (client *XXXArrClient) GetQueue(ctx context.Context) (queue int, err error) {
	var resp []byte
	resp, err = client.call(ctx, "/api/v3/queue")

	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		var stats struct{ TotalRecords int }
		if err = decoder.Decode(&stats); err == nil {
			queue = stats.TotalRecords
		}
	}

	log.WithError(err).WithFields(log.Fields{
		"application": client.Application,
		"queue":       queue,
	}).Debug("xxxarr GetQueue")

	return
}

// GetMonitored retrieves how many moves/series are being monitored & unmonitored
func (client *XXXArrClient) GetMonitored(ctx context.Context) (monitored int, unmonitored int, err error) {
	var endpoint string
	if client.Application == "sonarr" {
		endpoint = "/api/v3/series"
	} else if client.Application == "radarr" {
		endpoint = "/api/v3/movie"
	} else {
		panic("invalid application: " + client.Application)
	}

	var resp []byte
	resp, err = client.call(ctx, endpoint)

	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		var stats []struct{ Monitored bool }
		if err = decoder.Decode(&stats); err == nil {
			monitored = 0
			unmonitored = 0
			for _, stat := range stats {
				if stat.Monitored {
					monitored++
				} else {
					unmonitored++
				}

			}
		}
	}

	log.WithError(err).WithFields(log.Fields{
		"application": client.Application,
		"monitored":   monitored,
		"unmonitored": unmonitored,
	}).Debug("xxxarr GetMonitored")

	return
}

// call the specified Sonarr/Radarr API endpoint
func (client *XXXArrClient) call(ctx context.Context, endpoint string) (body []byte, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+endpoint, nil)
	req.Header.Add("X-Api-Key", client.APIKey)

	var timer *prometheus.Timer
	if client.Options.PrometheusSummary != nil {
		timer = prometheus.NewTimer(client.Options.PrometheusSummary.WithLabelValues(client.Application, endpoint))
	}

	var resp *http.Response
	resp, err = client.Client.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}

	if err == nil {
		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
		} else {
			err = fmt.Errorf("%s", resp.Status)
		}
		_ = resp.Body.Close()
	}
	return body, err
}
