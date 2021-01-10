package transmission_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/gotools/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"mediamon/internal/transmission"
)

func TestProbe_Run(t *testing.T) {
	probe := transmission.NewProbe("")
	probe.Client.Client = httpstub.NewTestClient(loopback)

	log.SetLevel(log.DebugLevel)

	probe.Run()

	value, _ := metrics.LoadValue("mediaserver_server_info", "transmission", "2.94 (d8e60ee44f)")
	assert.Equal(t, float64(1), value)
	value, _ = metrics.LoadValue("mediaserver_active_torrent_count")
	assert.Equal(t, float64(1), value)
	value, _ = metrics.LoadValue("mediaserver_paused_torrent_count")
	assert.Equal(t, float64(2), value)
	value, _ = metrics.LoadValue("mediaserver_download_speed")
	assert.Equal(t, float64(100), value)
	value, _ = metrics.LoadValue("mediaserver_upload_speed")
	assert.Equal(t, float64(25), value)
}

func TestFailingServer(t *testing.T) {
	probe := transmission.NewProbe("")
	probe.Client.Client = httpstub.NewTestClient(httpstub.Failing)

	assert.NotPanics(t, func() { probe.Run() })
}

func TestAuthentication(t *testing.T) {
	probe := transmission.NewProbe("")
	probe.Client.Client = httpstub.NewTestClient(loopback)

	log.SetLevel(log.DebugLevel)

	err := probe.Run()
	assert.Nil(t, err)

	// simulate the session key expiring
	probe.SessionID = "4321"
	err = probe.Run()
	assert.Nil(t, err)
	assert.Equal(t, "1234", probe.SessionID)
}

// Server loopback function

func loopback(req *http.Request) *http.Response {
	const sessionID = "1234"
	header := make(http.Header)
	header.Set("X-Transmission-Session-Id", sessionID)

	if req.Header.Get("X-Transmission-Session-Id") != sessionID {
		return &http.Response{
			StatusCode: 409,
			Status:     "No Session ID",
			Header:     header,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return &http.Response{
			StatusCode: 500,
			Status:     err.Error(),
			Header:     header,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}

	defer req.Body.Close()

	if string(body) == `{ "method": "session-get" }` {
		return &http.Response{
			StatusCode: 200,
			Header:     header,
			Body:       ioutil.NopCloser(bytes.NewBufferString(sessionGet)),
		}
	} else if string(body) == `{ "method": "session-stats" }` {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(sessionStats)),
		}
	} else {
		return &http.Response{
			StatusCode: 404,
			Status:     "Invalid method",
		}
	}
}

// Responses

const (
	sessionStats = `{
  "arguments": {
    "activeTorrentCount": 1,
    "cumulative-stats": {
      "downloadedBytes": 854551593612,
      "filesAdded": 1299,
      "secondsActive": 9650384,
      "sessionCount": 61,
      "uploadedBytes": 119748312546
    },
    "current-stats": {
      "downloadedBytes": 6435709785,
      "filesAdded": 20,
      "secondsActive": 589217,
      "sessionCount": 1,
      "uploadedBytes": 329128740
    },
    "downloadSpeed": 100,
    "pausedTorrentCount": 2,
    "torrentCount": 0,
    "uploadSpeed": 25
  },
  "result": "success"
}`

	sessionGet = `{
  "arguments": {
    "alt-speed-down": 50,
    "alt-speed-enabled": false,
    "alt-speed-time-begin": 540,
    "alt-speed-time-day": 127,
    "alt-speed-time-enabled": false,
    "alt-speed-time-end": 1020,
    "alt-speed-up": 50,
    "blocklist-enabled": false,
    "blocklist-size": 0,
    "blocklist-url": "http://www.example.com/blocklist",
    "cache-size-mb": 4,
    "config-dir": "/config",
    "dht-enabled": true,
    "download-dir": "/data/completed",
    "download-dir-free-space": 2466986983424,
    "download-queue-enabled": true,
    "download-queue-size": 5,
    "encryption": "preferred",
    "idle-seeding-limit": 5,
    "idle-seeding-limit-enabled": true,
    "incomplete-dir": "/data/incomplete",
    "incomplete-dir-enabled": true,
    "lpd-enabled": false,
    "peer-limit-global": 200,
    "peer-limit-per-torrent": 50,
    "peer-port": 51413,
    "peer-port-random-on-start": false,
    "pex-enabled": true,
    "port-forwarding-enabled": false,
    "queue-stalled-enabled": true,
    "queue-stalled-minutes": 30,
    "rename-partial-files": true,
    "rpc-version": 15,
    "rpc-version-minimum": 1,
    "script-torrent-done-enabled": false,
    "script-torrent-done-filename": "",
    "seed-queue-enabled": false,
    "seed-queue-size": 10,
    "seedRatioLimit": 0.0099,
    "seedRatioLimited": true,
    "speed-limit-down": 100,
    "speed-limit-down-enabled": false,
    "speed-limit-up": 100,
    "speed-limit-up-enabled": false,
    "start-added-torrents": true,
    "trash-original-torrent-files": false,
    "units": {
      "memory-bytes": 1024,
      "memory-units": [
        "KiB",
        "MiB",
        "GiB",
        "TiB"
      ],
      "size-bytes": 1000,
      "size-units": [
        "kB",
        "MB",
        "GB",
        "TB"
      ],
      "speed-bytes": 1000,
      "speed-units": [
        "kB/s",
        "MB/s",
        "GB/s",
        "TB/s"
      ]
    },
    "utp-enabled": false,
    "version": "2.94 (d8e60ee44f)"
  },
  "result": "success"
}`
)
