package xxxarr

import (
	"context"
	"github.com/clambin/mediamon/internal/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient"
	log "github.com/sirupsen/logrus"
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
func (probe *Probe) Run(ctx context.Context) error {
	var (
		err         error
		version     string
		count       int
		monitored   int
		unmonitored int
	)

	app := probe.GetApplication(ctx)

	probeLogger := log.WithField("application", app)

	// Get the version
	if version, err = probe.GetVersion(ctx); err != nil {
		probeLogger.WithField("err", err).Warning("could not get version")
	} else {
		metrics.MediaServerVersion.WithLabelValues(app, version).Set(1)
	}

	// Get the calendar
	if count, err = probe.GetCalendar(ctx); err != nil {
		probeLogger.WithField("err", err).Warning("could not get calendar")
	} else {
		metrics.XXXArrCalendarCount.WithLabelValues(app).Set(float64(count))
	}

	// Get queued series / movies
	if count, err = probe.GetQueue(ctx); err != nil {
		probeLogger.WithField("err", err).Warning("could not get queue")
	} else {
		metrics.XXXArrQueuedCount.WithLabelValues(app).Set(float64(count))
	}

	// Get monitored/unmonitored series / movies
	if monitored, unmonitored, err = probe.GetMonitored(ctx); err != nil {
		probeLogger.WithField("err", err).Warning("could not get monitored series/movies")
	} else {
		metrics.XXXArrMonitoredCount.WithLabelValues(app).Set(float64(monitored))
		metrics.XXXArrUnmonitoredCount.WithLabelValues(app).Set(float64(unmonitored))
	}

	return err
}
