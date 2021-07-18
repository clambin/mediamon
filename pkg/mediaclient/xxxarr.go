package mediaclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	GetApplication(ctx context.Context) string
}

// XXXArrClient calls the Sonarr/Radarr APIs.  Application specifies whether this is a
// Sonarr ("sonarr") or Radarr ("radarr") server.  XXXArrClient will panic if Application
// contains any other values.
type XXXArrClient struct {
	Client      *http.Client
	URL         string
	APIKey      string
	Application string
}

// GetApplication returns the client's configured application
func (client *XXXArrClient) GetApplication(_ context.Context) string {
	return client.Application
}

// GetVersion retrieves the version of the Sonarr/Radarr server
func (client *XXXArrClient) GetVersion(ctx context.Context) (string, error) {
	var (
		err   error
		resp  []byte
		stats struct {
			Version string
		}
	)

	if resp, err = client.call(ctx, "/api/v3/system/status"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{
		"err":         err,
		"application": client.Application,
		"version":     stats.Version,
	}).Debug("xxxarr Version")

	return stats.Version, err
}

// GetCalendar retrieves the number of upcoming movies/series airing today and tomorrow
func (client *XXXArrClient) GetCalendar(ctx context.Context) (int, error) {
	var (
		err      error
		resp     []byte
		calendar int
	)
	// TODO: add start/end date optional parameters?
	if resp, err = client.call(ctx, "/api/v3/calendar"); err == nil {
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

	log.WithFields(log.Fields{
		"err":         err,
		"application": client.Application,
		"calendar":    calendar,
	}).Debug("xxxarr getCalendar")

	return calendar, err
}

// GetQueue retrieves how many movies/series are currently downloading
func (client *XXXArrClient) GetQueue(ctx context.Context) (int, error) {
	var (
		err   error
		resp  []byte
		queue int
	)
	if resp, err = client.call(ctx, "/api/v3/queue"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		var stats struct{ TotalRecords int }
		if err = decoder.Decode(&stats); err == nil {
			queue = stats.TotalRecords
		}
	}

	log.WithFields(log.Fields{
		"err":         err,
		"application": client.Application,
		"queue":       queue,
	}).Debug("xxxarr GetQueue")

	return queue, err
}

// GetMonitored retrieves how many moves/series are being monitored & unmonitored
func (client *XXXArrClient) GetMonitored(ctx context.Context) (int, int, error) {
	var (
		err         error
		resp        []byte
		monitored   int
		unmonitored int
		endpoint    string
	)

	if client.Application == "sonarr" {
		endpoint = "/api/v3/series"
	} else if client.Application == "radarr" {
		endpoint = "/api/v3/movie"
	} else {
		panic("invalid application: " + client.Application)
	}

	if resp, err = client.call(ctx, endpoint); err == nil {
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

	log.WithFields(log.Fields{
		"err":         err,
		"application": client.Application,
		"monitored":   monitored,
		"unmonitored": unmonitored,
	}).Debug("xxxarr GetMonitored")

	return monitored, unmonitored, err
}

// call the specified Sonarr/Radarr API endpoint
func (client *XXXArrClient) call(ctx context.Context, endpoint string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+endpoint, nil)
	req.Header.Add("X-Api-Key", client.APIKey)

	if resp, err = client.Client.Do(req); err == nil {
		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
		} else {
			err = fmt.Errorf("%s", resp.Status)
		}
		_ = resp.Body.Close()
	}
	return body, err
}
