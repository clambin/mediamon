package mediaclient_test

import (
	"bytes"
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/mediamon/pkg/mediaclient"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestTransmissionClient_GetVersion(t *testing.T) {
	client := &mediaclient.TransmissionClient{Client: httpstub.NewTestClient(transmissionLoopback)}

	version, err := client.GetVersion()
	assert.Nil(t, err)
	assert.Equal(t, "2.94 (d8e60ee44f)", version)
}

func TestTransmissionClient_GetStats(t *testing.T) {
	client := &mediaclient.TransmissionClient{Client: httpstub.NewTestClient(transmissionLoopback)}

	active, paused, download, upload, err := client.GetStats()
	assert.Nil(t, err)
	assert.Equal(t, 1, active)
	assert.Equal(t, 2, paused)
	assert.Equal(t, 100, download)
	assert.Equal(t, 25, upload)
}

func TestTransmissionClient_Authentication(t *testing.T) {
	var (
		err        error
		oldVersion string
		newVersion string
	)

	client := &mediaclient.TransmissionClient{Client: httpstub.NewTestClient(transmissionLoopback)}
	//log.SetLevel(log.DebugLevel)

	oldVersion, err = client.GetVersion()
	assert.Nil(t, err)

	// simulate the session key expiring
	client.SessionID = "4321"

	newVersion, err = client.GetVersion()
	// call succeeded
	assert.Nil(t, err)
	// and the next SessionID has been set
	assert.Equal(t, "1234", client.SessionID)
	// and the call worked
	assert.Equal(t, oldVersion, newVersion)
}

func TestCallFailure(t *testing.T) {
	client := &mediaclient.TransmissionClient{Client: httpstub.NewTestClient(serverUnavailable)}

	assert.Eventually(t, func() bool {
		_, err := client.GetVersion()
		return err != nil

	}, 1*time.Second, 10*time.Millisecond)
}

// Server loopback function
func serverUnavailable(_ *http.Request) *http.Response {
	return nil
}

// Server loopback function
func transmissionLoopback(req *http.Request) *http.Response {
	const sessionID = "1234"
	header := make(http.Header)
	header.Set("X-Transmission-Session-Id", sessionID)

	if req.Header.Get("X-Transmission-Session-Id") != sessionID {
		return &http.Response{
			StatusCode: 409,
			Status:     "no/wrong Session ID",
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
