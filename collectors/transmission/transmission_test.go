package transmission_test

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/tests"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := transmission.NewCollector("http://localhost:8888", 5*time.Minute)
	metrics := make(chan *prometheus.Desc)
	go c.Describe(metrics)

	for _, metricName := range []string{
		"mediamon_transmission_version",
		"mediamon_transmission_active_torrent_count",
		"mediamon_transmission_paused_torrent_count",
		"mediamon_transmission_download_speed",
		"mediamon_transmission_upload_speed",
	} {
		metric := <-metrics
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	c := transmission.NewCollector("", time.Minute)
	c.(*transmission.Collector).TransmissionAPI = &server{}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	metric := <-metrics
	assert.True(t, tests.ValidateMetric(metric, 1, "version", "foo"))

	for _, value := range []float64{1, 2, 100, 25} {
		metric = <-metrics
		assert.True(t, tests.ValidateMetric(metric, value, "", ""))

	}
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
