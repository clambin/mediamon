package plex_test

import (
	"context"
	"errors"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/tests"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := plex.NewCollector("http://localhost:8888", "username", "password", 5*time.Minute)
	metrics := make(chan *prometheus.Desc)
	go c.Describe(metrics)

	for _, metricName := range []string{
		"mediamon_plex_version",
		"mediamon_plex_transcoder_encoding_count",
		"mediamon_plex_transcoder_speed_total",
		"mediamon_plex_session_count",
		"mediamon_plex_transcoder_type_count",
	} {
		metric := <-metrics
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	c := plex.NewCollector("", "", "", time.Minute)
	c.(*plex.Collector).PlexAPI = &client{}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	version := <-metrics
	assert.True(t, tests.ValidateMetric(version, 1, "version", "foo"))

	transcoding := <-metrics
	assert.True(t, tests.ValidateMetric(transcoding, 2, "", ""))

	speed := <-metrics
	assert.True(t, tests.ValidateMetric(speed, 3.1, "", ""))

	success := 0
	for i := 0; i < 2; i++ {
		user := <-metrics
		for _, name := range []string{"foo", "bar"} {
			if tests.ValidateMetric(user, 1, "user", name) {
				success++
			}
		}
	}
	assert.Equal(t, 2, success)

	success = 0
	for i := 0; i < 2; i++ {
		mode := <-metrics
		for _, name := range []string{"direct", "copy"} {
			if tests.ValidateMetric(mode, 1, "mode", name) {
				success++
			}
		}
	}
	assert.Equal(t, 2, success)
}

func TestCollector_Collect_Fail(t *testing.T) {
	c := plex.NewCollector("", "", "", time.Minute)
	c.(*plex.Collector).PlexAPI = &client{failing: true}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	assert.Never(t, func() bool { return len(metrics) > 0 }, 100*time.Millisecond, 10*time.Millisecond)

}

type client struct {
	failing bool
}

func (client *client) GetVersion(_ context.Context) (string, error) {
	if client.failing {
		return "", errors.New("failing")
	}
	return "foo", nil
}

func (client *client) GetSessions(_ context.Context) (users map[string]int, modes map[string]int, transcoding int, speed float64, err error) {
	if client.failing {
		return nil, nil, 0, 0.0, errors.New("failing")
	}
	users = make(map[string]int)
	modes = make(map[string]int)

	users["foo"] = 1
	users["bar"] = 1
	modes["direct"] = 1
	modes["copy"] = 1
	transcoding = 2
	speed = 3.1
	return
}
