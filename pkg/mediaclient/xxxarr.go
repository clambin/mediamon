package mediaclient

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// XXXArrAPI interface
type XXXArrAPI interface {
	GetVersion() (string, error)
	GetCalendar() (int, error)
	GetQueue() (int, error)
	GetMonitored() (int, int, error)
	GetApplication() string
}

// XXXArrClient to call the Sonarr/Radarr APIs
type XXXArrClient struct {
	Client      *http.Client
	URL         string
	APIKey      string
	Application string
}

// GetApplication returns the client's configured application
func (client *XXXArrClient) GetApplication() string {
	return client.Application
}

func (client *XXXArrClient) GetVersion() (string, error) {
	var (
		err   error
		resp  []byte
		stats struct {
			Version string
		}
	)

	if resp, err = client.call("/api/system/status"); err == nil {
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

func (client *XXXArrClient) GetCalendar() (int, error) {
	var (
		err      error
		resp     []byte
		calendar int
	)
	if resp, err = client.call("/api/calendar"); err == nil {
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

func (client *XXXArrClient) GetQueue() (int, error) {
	var (
		err   error
		resp  []byte
		queue int
	)
	if resp, err = client.call("/api/queue"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		var stats []struct{ Status string }
		if err = decoder.Decode(&stats); err == nil {
			queue = len(stats)
		}
	}

	log.WithFields(log.Fields{
		"err":         err,
		"application": client.Application,
		"queue":       queue,
	}).Debug("xxxarr GetQueue")

	return queue, err
}

func (client *XXXArrClient) GetMonitored() (int, int, error) {
	var (
		err         error
		resp        []byte
		monitored   int
		unmonitored int
		endpoint    string
	)

	if client.Application == "sonarr" {
		endpoint = "/api/series"
	} else if client.Application == "radarr" {
		endpoint = "/api/movie"
	} else {
		panic("invalid application: " + client.Application)
	}

	if resp, err = client.call(endpoint); err == nil {
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
func (client *XXXArrClient) call(endpoint string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	req, _ := http.NewRequest("GET", client.URL+endpoint, nil)
	req.Header.Add("X-Api-Key", client.APIKey)

	if resp, err = client.Client.Do(req); err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
		} else {
			err = errors.New(resp.Status)
		}
	}

	return body, err
}
