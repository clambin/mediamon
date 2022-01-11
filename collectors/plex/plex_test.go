package plex_test

import (
	"context"
	"errors"
	"github.com/clambin/mediamon/collectors/plex"
	plexAPI "github.com/clambin/mediamon/pkg/mediaclient/plex"
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
		"mediamon_plex_transcoder_total_count",
		"mediamon_plex_transcoder_active_count",
		"mediamon_plex_transcoder_speed_total",
		"mediamon_plex_session_location_count",
		"mediamon_plex_session_user_count",
	} {
		metric := <-metrics
		assert.Contains(t, metric.String(), "\""+metricName+"\"", metricName)
	}
}

func TestCollector_Collect(t *testing.T) {
	c := plex.NewCollector("", "", "", time.Minute)
	c.(*plex.Collector).API = &client{}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	version := <-metrics
	assert.True(t, tests.ValidateMetric(version, 1, "version", "foo"))

	transcoders := <-metrics
	assert.True(t, tests.ValidateMetric(transcoders, 2, "", ""))

	transcoding := <-metrics
	assert.True(t, tests.ValidateMetric(transcoding, 1, "", ""))

	speed := <-metrics
	assert.True(t, tests.ValidateMetric(speed, 21.0, "", ""))

	location := <-metrics
	assert.True(t, tests.ValidateMetric(location, 1.0, "location", "lan"))
	location = <-metrics
	assert.True(t, tests.ValidateMetric(location, 2.0, "location", "wan"))

	success := 0
	for i := 0; i < 2; i++ {
		user := <-metrics
		if tests.ValidateMetric(user, 2, "user", "foo") ||
			tests.ValidateMetric(user, 1, "user", "bar") {
			success++
		}
	}

	assert.Equal(t, 2, success)
}

func TestCollector_Collect_Fail(t *testing.T) {
	c := plex.NewCollector("", "", "", time.Minute)
	c.(*plex.Collector).API = &client{failing: true}

	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	metric := <-metrics
	assert.Equal(t, `Desc{fqName: "mediamon_error", help: "Error getting Plex version", constLabels: {}, variableLabels: []}`, metric.Desc().String())
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

func (client *client) GetSessions(_ context.Context) (sessions []plexAPI.Session, err error) {
	if client.failing {
		return nil, errors.New("failing")
	}

	sessions = []plexAPI.Session{
		{
			User:      "foo",
			Local:     true,
			Transcode: false,
		},
		{
			User:      "bar",
			Local:     false,
			Transcode: true,
			Throttled: false,
			Speed:     19.2,
		},
		{
			User:      "foo",
			Local:     false,
			Transcode: true,
			Throttled: true,
			Speed:     1.8,
		},
	}
	return
}
