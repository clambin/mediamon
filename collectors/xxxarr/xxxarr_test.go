package xxxarr_test

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/tests"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := xxxarr.NewCollector("http://localhost:8888", "", "sonarr", 5*time.Minute)
	metrics := make(chan *prometheus.Desc)
	go c.Describe(metrics)

	for _, metricName := range []string{
		"mediamon_xxxarr_version",
		"mediamon_xxxarr_calendar_count",
		"mediamon_xxxarr_queued_count",
		"mediamon_xxxarr_monitored_count",
		"mediamon_xxxarr_unmonitored_count",
	} {
		metric := <-metrics
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect_Sonarr(t *testing.T) {
	c := xxxarr.NewCollector("", "", "sonarr", 5*time.Minute)
	c.(*xxxarr.Collector).XXXArrAPI = &server{application: "sonarr"}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	metric := <-metrics
	assert.True(t, tests.ValidateMetric(metric, 1, "version", "foo"))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 5, "", ""))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 2, "", ""))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 10, "", ""))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 3, "", ""))
}

func TestCollector_Collect_Radarr(t *testing.T) {
	c := xxxarr.NewCollector("", "", "radarr", 5*time.Minute)
	c.(*xxxarr.Collector).XXXArrAPI = &server{application: "radarr"}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	metric := <-metrics
	assert.True(t, tests.ValidateMetric(metric, 1, "version", "foo"))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 5, "", ""))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 2, "", ""))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 10, "", ""))

	metric = <-metrics
	assert.True(t, tests.ValidateMetric(metric, 3, "", ""))
}

func TestCollector_Bad_Application(t *testing.T) {
	assert.Panics(t, func() {
		_ = xxxarr.NewCollector("", "", "foo", 5*time.Minute)
	})
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

func (s *server) GetApplication(_ context.Context) (application string) {
	return s.application
}
