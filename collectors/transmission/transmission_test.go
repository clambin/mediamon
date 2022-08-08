package transmission_test

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/collectors/transmission"
	transmission2 "github.com/clambin/mediamon/pkg/mediaclient/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := transmission.NewCollector("http://localhost:8888")
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
	c := transmission.NewCollector("")
	c.(*transmission.Collector).API = &server{}

	e := strings.NewReader(`# HELP mediamon_transmission_active_torrent_count Number of active torrents
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

func TestCollector_Collect_Fail(t *testing.T) {
	c := transmission.NewCollector("")
	c.(*transmission.Collector).API = &server{fail: true}

	err := testutil.CollectAndCompare(c, strings.NewReader(``))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `Desc{fqName: "mediamon_error", help: "Error getting transmission metrics", constLabels: {}, variableLabels: []}`)
}

type server struct {
	fail bool
}

func (server server) GetSessionParameters(_ context.Context) (response transmission2.SessionParameters, err error) {
	if server.fail {
		err = fmt.Errorf("failed")
		return
	}
	response.Arguments.Version = "foo"
	return
}

func (server server) GetSessionStatistics(_ context.Context) (stats transmission2.SessionStats, err error) {
	if server.fail {
		err = fmt.Errorf("failed")
		return
	}
	stats.Arguments.ActiveTorrentCount = 1
	stats.Arguments.PausedTorrentCount = 2
	stats.Arguments.UploadSpeed = 25
	stats.Arguments.DownloadSpeed = 100
	return
}
