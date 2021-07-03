package transmission

import (
	"context"
	"github.com/clambin/mediamon/internal/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Probe to measure Transmission metrics
type Probe struct {
	mediaclient.TransmissionAPI
}

// NewProbe creates a new Probe
func NewProbe(url string) *Probe {
	return &Probe{&mediaclient.TransmissionClient{Client: &http.Client{}, URL: url}}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run(ctx context.Context) (err error) {
	var (
		version        string
		activeTorrents int
		pausedTorrents int
		downloadSpeed  int
		uploadSpeed    int
	)

	// Get the version
	if version, err = probe.GetVersion(ctx); err == nil {
		metrics.MediaServerVersion.WithLabelValues("transmission", version).Set(1)

		// Get statistics
		activeTorrents, pausedTorrents, downloadSpeed, uploadSpeed, err = probe.GetStats(ctx)
	}

	if err == nil {
		metrics.TransmissionActiveTorrentCount.Set(float64(activeTorrents))
		metrics.TransmissionPausedTorrentCount.Set(float64(pausedTorrents))
		metrics.TransmissionDownloadSpeed.Set(float64(downloadSpeed))
		metrics.TransmissionUploadSpeed.Set(float64(uploadSpeed))
	}

	if err != nil {
		log.WithField("err", err).Warning("Could not get Transmission version")
	}
	return
}
