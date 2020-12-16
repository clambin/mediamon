package transmission

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"mediamon/internal/metrics"
	"net/http"
)

type Probe struct {
	client *Client
}

func NewProbe(url string) *Probe {
	return NewProbeWithHTTPClient(&http.Client{}, url)
}

func NewProbeWithHTTPClient(client *http.Client, url string) *Probe {
	return &Probe{client: NewAPIWithHTTPClient(client, url)}
}

func (probe *Probe) Run() {
	// Get the version
	if version, err := probe.getVersion(); err != nil {
		log.Warningf("Could not get Transmission version: %s", err)
	} else {
		metrics.Publish("version", 1, "transmission", version)
	}

	if activeTorrents, pausedTorrents, downloadSpeed, uploadSpeed, err := probe.getStats(); err != nil {
		log.Warningf("Could not get Transmission Statistics: %s", err)
	} else {
		metrics.Publish("active_torrent_count", float64(activeTorrents), "transmission")
		metrics.Publish("paused_torrent_count", float64(pausedTorrents), "transmission")
		metrics.Publish("download_speed", float64(downloadSpeed), "transmission")
		metrics.Publish("upload_speed", float64(uploadSpeed), "transmission")
	}
}

func (probe *Probe) getVersion() (string, error) {
	var stats = struct {
		Arguments struct {
			Version string
		}
		Result string
	}{}

	resp, err := probe.client.Call("session-get")
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

	resp, err := probe.client.Call("session-stats")
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
