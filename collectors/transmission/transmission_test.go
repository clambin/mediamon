package transmission_test

import (
	"context"
	"fmt"
	"github.com/clambin/go-metrics/tools"
	"github.com/clambin/mediamon/collectors/transmission"
	transmission2 "github.com/clambin/mediamon/pkg/mediaclient/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
	assert.Equal(t, 1.0, tools.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, "foo", tools.MetricLabel(metric, "version"))

	for _, value := range []float64{1, 2, 100, 25} {
		metric = <-ch
		assert.Equal(t, value, tools.MetricValue(metric).GetGauge().GetValue())

	}
}

func TestCollector_Collect_Fail(t *testing.T) {
	c := transmission.NewCollector("")
	c.(*transmission.Collector).API = &server{fail: true}

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	assert.Never(t, func() bool {
		<-ch
		return true
	}, 100*time.Millisecond, 10*time.Millisecond)
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
