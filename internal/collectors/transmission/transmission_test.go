package transmission_test

import (
	"github.com/clambin/mediamon/v2/internal/collectors/transmission"
	"github.com/clambin/mediamon/v2/internal/collectors/transmission/mocks"
	"github.com/hekmon/transmissionrpc/v3"
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
	g := mocks.NewTransmissionClient(t)
	sessionStats := transmissionrpc.SessionStats{
		ActiveTorrentCount: 1,
		PausedTorrentCount: 2,
		UploadSpeed:        25,
		DownloadSpeed:      100,
	}
	g.EXPECT().SessionStats(mock.Anything).Return(sessionStats, nil)
	sessionArguments := transmissionrpc.SessionArguments{Version: constP("foo")}
	g.EXPECT().SessionArgumentsGetAll(mock.Anything).Return(sessionArguments, nil)

	c := transmission.NewCollector("", slog.Default())
	c.Collector.(*transmission.Collector).TransmissionClient = g

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

func constP[T any](t T) *T {
	return &t
}
