package transmission_test

import (
	transmissionClient "github.com/clambin/mediaclients/transmission"
	"github.com/clambin/mediamon/v2/internal/collectors/transmission"
	"github.com/clambin/mediamon/v2/internal/collectors/transmission/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"strings"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := transmission.NewCollector("http://localhost:8888", slog.Default())
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{
		"mediamon_transmission_version",
		"mediamon_transmission_active_torrent_count",
		"mediamon_transmission_paused_torrent_count",
		"mediamon_transmission_download_speed",
		"mediamon_transmission_upload_speed",
	} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	g := mocks.NewGetter(t)
	var sessionStats transmissionClient.SessionStats
	sessionStats.Arguments.ActiveTorrentCount = 1
	sessionStats.Arguments.PausedTorrentCount = 2
	sessionStats.Arguments.UploadSpeed = 25
	sessionStats.Arguments.DownloadSpeed = 100
	g.EXPECT().GetSessionStatistics(mock.Anything).Return(sessionStats, nil)
	var sessionParameters transmissionClient.SessionParameters
	sessionParameters.Arguments.Version = "foo"
	g.EXPECT().GetSessionParameters(mock.Anything).Return(sessionParameters, nil)

	c := transmission.NewCollector("", slog.Default())
	c.Collector.(*transmission.Collector).Transmission = g

	e := strings.NewReader(`
# HELP circuit_breaker_consecutive_errors consecutive errors
# TYPE circuit_breaker_consecutive_errors gauge
circuit_breaker_consecutive_errors{circuit_breaker="transmission"} 0

# HELP circuit_breaker_consecutive_successes consecutive successes
# TYPE circuit_breaker_consecutive_successes gauge
circuit_breaker_consecutive_successes{circuit_breaker="transmission"} 1

# HELP circuit_breaker_state state of the circuit breaker (0: closed, 1:open, 2:half-open)
# TYPE circuit_breaker_state gauge
circuit_breaker_state{circuit_breaker="transmission"} 0

# HELP mediamon_transmission_active_torrent_count Number of active torrents
# TYPE mediamon_transmission_active_torrent_count gauge
mediamon_transmission_active_torrent_count{url=""} 1

# HELP mediamon_transmission_download_speed Transmission download speed in bytes / sec
# TYPE mediamon_transmission_download_speed gauge
mediamon_transmission_download_speed{url=""} 100

# HELP mediamon_transmission_paused_torrent_count Number of paused torrents
# TYPE mediamon_transmission_paused_torrent_count gauge
mediamon_transmission_paused_torrent_count{url=""} 2

# HELP mediamon_transmission_upload_speed Transmission upload speed in bytes / sec
# TYPE mediamon_transmission_upload_speed gauge
mediamon_transmission_upload_speed{url=""} 25

# HELP mediamon_transmission_version version info
# TYPE mediamon_transmission_version gauge
mediamon_transmission_version{url="",version="foo"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(c, e))
}
