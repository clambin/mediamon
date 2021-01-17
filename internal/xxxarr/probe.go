package xxxarr

import (
	log "github.com/sirupsen/logrus"
	"mediamon/internal/metrics"
	"mediamon/pkg/mediaclient"
	"net/http"
)

// Probe to measure sonarr/radarr metrics
type Probe struct {
	mediaclient.XXXArrAPI
}

// NewProbe creates a new Probe
func NewProbe(url string, apiKey string, application string) *Probe {
	return &Probe{&mediaclient.XXXArrClient{Client: &http.Client{}, URL: url, APIKey: apiKey, Application: application}}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() error {
	var (
		err     error
		version string
		count   int
	)

	probeLogger := log.WithFields(log.Fields{"err": err, "application": probe.GetApplication()})

	// Get the version
	if version, err = probe.GetVersion(); err != nil {
		probeLogger.Warning("could not get version")
	} else {
		metrics.MediaServerVersion.WithLabelValues(probe.GetApplication(), version).Set(1)
	}

	// Get the calendar
	if count, err = probe.GetCalendar(); err != nil {
		probeLogger.Warning("could not get calendar")
	} else {
		metrics.XXXArrCalendarCount.WithLabelValues(probe.GetApplication()).Set(float64(count))
	}

	// Get queued series / movies
	if count, err = probe.GetQueue(); err != nil {
		probeLogger.Warning("could not get queue")
	} else {
		metrics.XXXArrQueuedCount.WithLabelValues(probe.GetApplication()).Set(float64(count))
	}

	// Get monitored/unmonitored series / movies
	if monitored, unmonitored, err := probe.GetMonitored(); err != nil {
		probeLogger.Warning("could not get monitored series/movies")
	} else {
		metrics.XXXArrMonitoredCount.WithLabelValues(probe.GetApplication()).Set(float64(monitored))
		metrics.XXXArrUnmonitoredCount.WithLabelValues(probe.GetApplication()).Set(float64(unmonitored))
	}

	return err
}
