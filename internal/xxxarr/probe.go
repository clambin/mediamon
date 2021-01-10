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
	var (
		err     error
		version string
		count   int
	)

	probeLogger := log.WithFields(log.Fields{"err": err, "application": probe.Application})

	// Get the version
	if version, err = probe.getVersion(); err != nil {
		probeLogger.Warning("could not get version")
	} else {
		metrics.MediaServerVersion.WithLabelValues(probe.Application, version).Set(1)
	}

	// Get the calendar
	if count, err = probe.getCalendar(); err != nil {
		probeLogger.Warning("could not get calendar")
	} else {
		metrics.XXXArrCalendarCount.WithLabelValues(probe.Application).Set(float64(count))
	}

	// Get queued series / movies
	if count, err = probe.getQueue(); err != nil {
		probeLogger.Warning("could not get queue")
	} else {
		metrics.XXXarrQueuedCount.WithLabelValues(probe.Application).Set(float64(count))
	}

	// Get monitored/unmonitored series / movies
	if monitored, unmonitored, err := probe.getMonitored(); err != nil {
		probeLogger.Warning("could not get monitored series/movies")
	} else {
		metrics.XXXarrMonitoredCount.WithLabelValues(probe.Application).Set(float64(monitored))
		metrics.XXXarrUnmonitoredCount.WithLabelValues(probe.Application).Set(float64(unmonitored))
	}
}

func (probe *Probe) getVersion() (string, error) {
	var (
		err   error
		resp  []byte
		stats struct {
			Version string
		}
	)

	if resp, err = probe.call("/api/system/status"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{
		"err":         err,
		"application": probe.Application,
		"version":     stats.Version,
	}).Debug("xxxarr Version")

	return stats.Version, err
}

func (probe *Probe) getCalendar() (int, error) {
	var (
		err      error
		resp     []byte
		calendar int
	)
	if resp, err = probe.call("/api/calendar"); err == nil {
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
		"application": probe.Application,
		"calendar":    calendar,
	}).Debug("xxxarr getCalendar")

	return calendar, err
}

func (probe *Probe) getQueue() (int, error) {
	var (
		err   error
		resp  []byte
		queue int
	)
	if resp, err = probe.call("/api/queue"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		var stats []struct{ Status string }
		if err = decoder.Decode(&stats); err == nil {
			queue = len(stats)
		}
	}

	log.WithFields(log.Fields{
		"err":         err,
		"application": probe.Application,
		"queue":       queue,
	}).Debug("xxxarr getQueue")

	return queue, err
}

func (probe *Probe) getMonitored() (int, int, error) {
	var (
		err         error
		resp        []byte
		monitored   int
		unmonitored int
	)

	endpoint := "/api/movie"
	if probe.Application == "sonarr" {
		endpoint = "/api/series"
	}

	if resp, err = probe.call(endpoint); err == nil {
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
		"application": probe.Application,
		"monitored":   monitored,
		"unmonitored": unmonitored,
	}).Debug("xxxarr getMonitored")

	return monitored, unmonitored, err
}
