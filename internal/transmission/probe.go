package transmission

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"mediamon/internal/metrics"
)

// Probe to measure Transmission metrics
type Probe struct {
	Client
}

// NewProbe creates a new Probe
func NewProbe(url string) *Probe {
	return &Probe{Client{Client: &http.Client{}, URL: url}}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() error {
	var (
		err            error
		version        string
		activeTorrents int
		pausedTorrents int
		downloadSpeed  int
		uploadSpeed    int
	)
	// Get the version
	if version, err = probe.getVersion(); err != nil {
		log.WithField("err", err).Warning("Could not get Transmission version")
	} else {
		metrics.MediaServerVersion.WithLabelValues("transmission", version).Set(1)
	}

	if activeTorrents, pausedTorrents, downloadSpeed, uploadSpeed, err = probe.getStats(); err != nil {
		log.WithField("err", err).Warning("Could not get Transmission Statistics")
	} else {
		metrics.TransmissionActiveTorrentCount.Set(float64(activeTorrents))
		metrics.TransmissionPausedTorrentCount.Set(float64(pausedTorrents))
		metrics.TransmissionDownloadSpeed.Set(float64(downloadSpeed))
		metrics.TransmissionUploadSpeed.Set(float64(uploadSpeed))
	}

	return err
}

func (probe *Probe) getVersion() (string, error) {
	var (
		err   error
		resp  []byte
		stats = struct {
			Arguments struct {
				Version string
			}
		}{}
	)

	if resp, err = probe.call("session-get"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{"err": err, "version": stats.Arguments.Version}).Debug("transmission getVersion")

	return stats.Arguments.Version, err
}

func (probe *Probe) getStats() (int, int, int, int, error) {
	var (
		err   error
		resp  []byte
		stats = struct {
			Arguments struct {
				ActiveTorrentCount int
				PausedTorrentCount int
				UploadSpeed        int
				DownloadSpeed      int
			}
			Result string
		}{}
	)

	if resp, err = probe.call("session-stats"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{"err": err, "stats": stats.Arguments}).Debug("transmission getStats")

	return stats.Arguments.ActiveTorrentCount,
		stats.Arguments.PausedTorrentCount,
		stats.Arguments.DownloadSpeed,
		stats.Arguments.UploadSpeed,
		err
}
