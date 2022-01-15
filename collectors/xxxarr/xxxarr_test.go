package xxxarr_test

import (
	"context"
	"fmt"
	"github.com/clambin/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testCollectorDescribe(t *testing.T, collector prometheus.Collector, labelString string) {
	ch := make(chan *prometheus.Desc)
	go collector.Describe(ch)

	for _, metricName := range []string{
		"mediamon_xxxarr_version",
		"mediamon_xxxarr_calendar_count",
		"mediamon_xxxarr_queued_count",
		"mediamon_xxxarr_monitored_count",
		"mediamon_xxxarr_unmonitored_count",
	} {
		metric := <-ch
		metricAsString := metric.String()
		assert.Contains(t, metricAsString, "\""+metricName+"\"")
		assert.Contains(t, metricAsString, labelString)
	}
}

func testCollectorCollect(t *testing.T, collector prometheus.Collector, application string) {
	ch := make(chan prometheus.Metric)
	go collector.Collect(ch)

	metric := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, "foo", metrics.MetricLabel(metric, "version"))

	for _, value := range []float64{5, 2, 10, 3} {
		metric = <-ch
		assert.Equal(t, value, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, application, metrics.MetricLabel(metric, "application"))
	}
}

type server struct {
	application string
	failing     bool
}

func (s *server) GetVersion(_ context.Context) (version string, err error) {
	if s.failing {
		return "", fmt.Errorf("failing")
	}
	return "foo", nil
}

func (s *server) GetCalendar(_ context.Context) (count int, err error) {
	if s.failing {
		return 0, fmt.Errorf("failing")
	}
	return 5, nil
}

func (s *server) GetQueue(_ context.Context) (count int, err error) {
	if s.failing {
		return 0, fmt.Errorf("failing")
	}
	return 2, nil
}

func (s *server) GetMonitored(_ context.Context) (monitored, unmonitored int, err error) {
	if s.failing {
		return 0, 0, fmt.Errorf("failing")
	}
	return 10, 3, nil
}

func (s *server) GetApplication() (application string) {
	return s.application
}

func (s *server) GetURL() (url string) {
	return "https://localhost:4321"
}
