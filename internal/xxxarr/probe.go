package xxxarr

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"mediamon/internal/metrics"
)

// Probe to measure sonarr/radarr metrics
type Probe struct {
	Client
	Application string
}

// NewProbe creates a new Probe
func NewProbe(url string, apiKey string, application string) *Probe {
	if isValid(application) == false {
		panic(application)
	}
	return &Probe{Client{Client: &http.Client{}, URL: url, APIKey: apiKey}, application}
}

func isValid(application string) bool {
	for _, valid := range []string{"sonarr", "radarr"} {
		if valid == application {
			return true
		}
	}
	return false
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() {
	// Get the version
	if version, err := probe.getVersion(); err != nil {
		log.Warningf("could not get %s version: %s", probe.Application, err)
	} else {
		metrics.MediaServerVersion.WithLabelValues(probe.Application, version).Set(1)
	}

	// Get the calendar
	if count, err := probe.getCalendar(); err != nil {
		log.Warningf("could not get %s calendar: %s", probe.Application, err)
	} else {
		metrics.XXXArrCalendarCount.WithLabelValues(probe.Application).Set(float64(count))
	}

	// Get queued series / movies
	if count, err := probe.getQueue(); err != nil {
		log.Warningf("could not get %s queue: %s", probe.Application, err)
	} else {
		metrics.XXXarrQueuedCount.WithLabelValues(probe.Application).Set(float64(count))
	}

	// Get monitored/unmonitored series / movies
	if monitored, unmonitored, err := probe.getMonitored(); err != nil {
		log.Warningf("could not get %s monitored series/movies: %s", probe.Application, err)
	} else {
		metrics.XXXarrMonitoredCount.WithLabelValues(probe.Application).Set(float64(monitored))
		metrics.XXXarrUnmonitoredCount.WithLabelValues(probe.Application).Set(float64(unmonitored))
	}
}

func (probe *Probe) getVersion() (string, error) {
	var stats = struct {
		Version string
	}{}

	resp, err := probe.call("/api/system/status")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			return stats.Version, nil
		}
	}
	return "", err
}

func (probe *Probe) getCalendar() (int, error) {
	var stats []struct {
		HasFile bool
	}

	resp, err := probe.call("/api/calendar")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			calendar := 0
			for _, stat := range stats {
				if stat.HasFile == false {
					calendar++
				}
			}
			return calendar, nil
		}
	}
	return 0, err
}

func (probe *Probe) getQueue() (int, error) {
	var stats []struct {
		Status string
	}

	resp, err := probe.call("/api/queue")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			return len(stats), nil
		}
	}
	return 0, err
}

func (probe *Probe) getMonitored() (int, int, error) {
	var stats []struct {
		Monitored bool
	}

	endpoint := "/api/movie"
	if probe.Application == "sonarr" {
		endpoint = "/api/series"
	}

	resp, err := probe.call(endpoint)
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			monitored := 0
			unmonitored := 0
			for _, stat := range stats {
				if stat.Monitored {
					monitored++
				} else {
					unmonitored++
				}

			}
			return monitored, unmonitored, nil
		}
	}
	return 0, 0, err
}
