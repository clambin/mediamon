package transmission_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/transmission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransmissionClient_GetSessionParameters(t *testing.T) {
	s := server{sessionID: "1234"}
	testServer := s.start()
	defer testServer.Close()

	c := transmission.NewClient(testServer.URL, nil)

	params, err := c.GetSessionParameters(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "2.94 (d8e60ee44f)", params.Arguments.Version)
}

func TestTransmissionClient_GetSessionStats(t *testing.T) {
	s := server{sessionID: "1234"}
	testServer := s.start()
	defer testServer.Close()

	c := transmission.NewClient(testServer.URL, nil)

	stats, err := c.GetSessionStatistics(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, stats.Arguments.ActiveTorrentCount)
	assert.Equal(t, 2, stats.Arguments.PausedTorrentCount)
	assert.Equal(t, 100, stats.Arguments.DownloadSpeed)
	assert.Equal(t, 25, stats.Arguments.UploadSpeed)
}

func TestTransmissionClient_Failures(t *testing.T) {
	s := server{sessionID: "1234", invalid: true}
	testServer := s.start()

	c := transmission.NewClient(testServer.URL, nil)

	ctx := context.Background()

	_, err := c.GetSessionParameters(ctx)
	assert.Error(t, err)

	s.invalid = false
	_, err = c.GetSessionParameters(ctx)
	assert.NoError(t, err)

	s.notSuccess = true
	_, err = c.GetSessionParameters(ctx)
	assert.Error(t, err)
	_, err = c.GetSessionStatistics(ctx)
	assert.Error(t, err)

	s.fail = true
	_, err = c.GetSessionParameters(ctx)
	assert.Error(t, err)

	testServer.Close()
	_, err = c.GetSessionParameters(ctx)
	assert.Error(t, err)
}

// Server handlers

type server struct {
	sessionID  string
	fail       bool
	invalid    bool
	notSuccess bool
}

func (s *server) start() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(s.handler))
}

func (s *server) handler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		_ = req.Body.Close()
	}()

	if s.fail {
		http.Error(w, "server broken", http.StatusInternalServerError)
		return
	}

	w.Header()["X-Transmission-Session-Id"] = []string{s.sessionID}
	if req.Header.Get("X-Transmission-Session-Id") != s.sessionID {
		http.Error(w, "Forbidden", http.StatusConflict)
		return
	}

	if s.invalid {
		_, _ = w.Write([]byte(`invalid content`))
		return
	}

	if s.notSuccess {
		_, _ = w.Write([]byte(`{ "result": "failed" }`))
		return
	}

	body, _ := io.ReadAll(req.Body)
	response, ok := transmissionResponses[string(body)]

	if ok == false {
		http.Error(w, "endpoint not implemented", http.StatusNotFound)
	} else {
		_, _ = w.Write([]byte(response))

	}
}

// Responses

var transmissionResponses = map[string]string{
	`{ "method": "session-get" }`: `{ "result": "success", "arguments": {
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
      		"memory-units": [ "KiB", "MiB", "GiB", "TiB" ],
      		"size-bytes": 1000,
      		"size-units": [ "kB", "MB", "GB", "TB" ],
      		"speed-bytes": 1000,
      		"speed-units": [ "kB/s", "MB/s", "GB/s", "TB/s" ]
		},
		"utp-enabled": false,
    	"version": "2.94 (d8e60ee44f)"
  	}}`,
	`{ "method": "session-stats" }`: `{ "result": "success", "arguments": {
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
	}}`,
}
