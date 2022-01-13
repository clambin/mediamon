package xxxarr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/metrics"
	"net/http"
)

// API interface
type API interface {
	GetVersion(ctx context.Context) (string, error)
	GetCalendar(ctx context.Context) (int, error)
	GetQueue(ctx context.Context) (int, error)
	GetMonitored(ctx context.Context) (int, int, error)
	GetApplication() string
	GetURL() string
}

// Client calls the Sonarr/Radarr APIs.  Application specifies whether this is a
// Sonarr ("sonarr") or Radarr ("radarr") server.  Client will panic if Application
// contains any other values.
type Client struct {
	Client      *http.Client
	URL         string
	APIKey      string
	Application string
	Options     Options
}

// Options contains options to alter Client behaviour
type Options struct {
	PrometheusMetrics metrics.PrometheusMetrics
}

// GetApplication returns the client's configured application
func (client *Client) GetApplication() string {
	return client.Application
}

// GetURL returns the server URL
func (client *Client) GetURL() string {
	return client.URL
}

// GetVersion retrieves the version of the Sonarr/Radarr server
func (client *Client) GetVersion(ctx context.Context) (version string, err error) {
	var stats struct {
		Version string
	}
	err = client.call(ctx, "/api/v3/system/status", &stats)

	if err == nil {
		version = stats.Version
	}
	return
}

// GetCalendar retrieves the number of upcoming movies/series airing today and tomorrow
func (client *Client) GetCalendar(ctx context.Context) (calendar int, err error) {
	var stats []struct {
		HasFile bool
	}
	err = client.call(ctx, "/api/v3/calendar", &stats)

	if err == nil {
		for _, stat := range stats {
			if stat.HasFile == false {
				calendar++
			}
		}
	}
	return
}

// GetQueue retrieves how many movies/series are currently downloading
func (client *Client) GetQueue(ctx context.Context) (queue int, err error) {
	var stats struct {
		TotalRecords int
	}
	err = client.call(ctx, "/api/v3/queue", &stats)

	if err == nil {
		queue = stats.TotalRecords
	}
	return
}

// GetMonitored retrieves how many moves/series are being monitored & unmonitored
func (client *Client) GetMonitored(ctx context.Context) (monitored int, unmonitored int, err error) {
	var endpoint string
	if client.Application == "sonarr" {
		endpoint = "/api/v3/series"
	} else if client.Application == "radarr" {
		endpoint = "/api/v3/movie"
	} else {
		panic("invalid application: " + client.Application)
	}

	var stats []struct {
		Monitored bool
	}
	err = client.call(ctx, endpoint, &stats)

	if err == nil {
		for _, stat := range stats {
			if stat.Monitored {
				monitored++
			} else {
				unmonitored++
			}
		}
	}
	return
}

// call the specified Sonarr/Radarr API endpoint
func (client *Client) call(ctx context.Context, endpoint string, response interface{}) (err error) {
	defer func() {
		client.Options.PrometheusMetrics.ReportErrors(err, client.Application, endpoint)
	}()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+endpoint, nil)
	req.Header.Add("X-Api-Key", client.APIKey)

	timer := client.Options.PrometheusMetrics.MakeLatencyTimer(client.Application, endpoint)

	var resp *http.Response
	resp, err = client.Client.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}

	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(response)
	_ = resp.Body.Close()

	return
}
