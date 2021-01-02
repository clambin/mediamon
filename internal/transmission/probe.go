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
func (probe *Probe) Run() {
	// Get the version
	if version, err := probe.getVersion(); err != nil {
		log.Warningf("Could not get Transmission version: %s", err)
	} else {
		metrics.MediaServerVersion.WithLabelValues("transmission", version).Set(1)
	}

	if activeTorrents, pausedTorrents, downloadSpeed, uploadSpeed, err := probe.getStats(); err != nil {
		log.Warningf("Could not get Transmission Statistics: %s", err)
	} else {
		metrics.TransmissionActiveTorrentCount.Set(float64(activeTorrents))
		metrics.TransmissionPausedTorrentCount.Set(float64(pausedTorrents))
		metrics.TransmissionDownloadSpeed.Set(float64(downloadSpeed))
		metrics.TransmissionUploadSpeed.Set(float64(uploadSpeed))
	}
}

func (probe *Probe) getVersion() (string, error) {
	var stats = struct {
		Arguments struct {
			Version string
		}
		Result string
	}{}

	resp, err := probe.call("session-get")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			return stats.Arguments.Version, nil
		}
	}
	return "", err
}

func (probe *Probe) getStats() (int, int, int, int, error) {
	var stats = struct {
		Arguments struct {
			ActiveTorrentCount int
			PausedTorrentCount int
			UploadSpeed        int
			DownloadSpeed      int
		}
		Result string
	}{}

	resp, err := probe.call("session-stats")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			return stats.Arguments.ActiveTorrentCount,
				stats.Arguments.PausedTorrentCount,
				stats.Arguments.DownloadSpeed,
				stats.Arguments.UploadSpeed,
				nil
		}
	}
	return 0, 0, 0, 0, err
}
