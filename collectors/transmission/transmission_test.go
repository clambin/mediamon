package transmission_test

import (
	"context"
	"fmt"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := transmission.NewCollector("http://localhost:8888", 5*time.Minute)
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
	c := transmission.NewCollector("", time.Minute)
	c.(*transmission.Collector).API = &server{}

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, "foo", metrics.MetricLabel(metric, "version"))

	for _, value := range []float64{1, 2, 100, 25} {
		metric = <-ch
		assert.Equal(t, value, metrics.MetricValue(metric).GetGauge().GetValue())

	}
}

func TestCollector_Collect_Fail(t *testing.T) {
	c := transmission.NewCollector("", time.Minute)
	c.(*transmission.Collector).API = &server{fail: true}

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	assert.Never(t, func() bool {
		_ = <-ch
		return true
	}, 100*time.Millisecond, 10*time.Millisecond)
}

type server struct {
	fail bool
}

func (server *server) GetVersion(_ context.Context) (string, error) {
	if server.fail {
		return "", fmt.Errorf("failed")
	}
	return "foo", nil
}

func (server *server) GetStats(_ context.Context) (int, int, int, int, error) {
	if server.fail {
		return 0, 0, 0, 0, fmt.Errorf("failed")
	}
	return 1, 2, 100, 25, nil
}
