package transmission

import (
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
func (probe *Probe) Run() (err error) {
	var (
		version        string
		activeTorrents int
		pausedTorrents int
		downloadSpeed  int
		uploadSpeed    int
	)

	// Get the version
	if version, err = probe.GetVersion(); err != nil {
		log.WithField("err", err).Warning("Could not get Transmission version")
	} else {
		metrics.MediaServerVersion.WithLabelValues("transmission", version).Set(1)

		// Get statistics
		if activeTorrents, pausedTorrents, downloadSpeed, uploadSpeed, err = probe.GetStats(); err != nil {
			log.WithField("err", err).Warning("Could not get Transmission Statistics")
		} else {
			metrics.TransmissionActiveTorrentCount.Set(float64(activeTorrents))
			metrics.TransmissionPausedTorrentCount.Set(float64(pausedTorrents))
			metrics.TransmissionDownloadSpeed.Set(float64(downloadSpeed))
			metrics.TransmissionUploadSpeed.Set(float64(uploadSpeed))
		}
	}
	return
}
