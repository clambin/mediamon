package plex_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/collectors/plex"
	plexAPI "github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := plex.NewCollector("http://localhost:8888", "username", "password", 5*time.Minute)
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{
		"mediamon_plex_version",
		"mediamon_plex_transcoder_total_count",
		"mediamon_plex_transcoder_active_count",
		"mediamon_plex_transcoder_speed_total",
		"mediamon_plex_session_location_count",
		"mediamon_plex_session_user_count",
	} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"", metricName)
	}
}

func TestCollector_Collect(t *testing.T) {
	c := plex.NewCollector("", "", "", time.Minute)
	c.(*plex.Collector).API = &client{}

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	version := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(version).GetGauge().GetValue())
	assert.Equal(t, "foo", metrics.MetricLabel(version, "version"))

	transcoders := <-ch
	assert.Equal(t, 2.0, metrics.MetricValue(transcoders).GetGauge().GetValue())

	transcoding := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(transcoding).GetGauge().GetValue())

	speed := <-ch
	assert.Equal(t, 21.0, metrics.MetricValue(speed).GetGauge().GetValue())

	location := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(location).GetGauge().GetValue())
	assert.Equal(t, "lan", metrics.MetricLabel(location, "location"))
	location = <-ch
	assert.Equal(t, 2.0, metrics.MetricValue(location).GetGauge().GetValue())
	assert.Equal(t, "wan", metrics.MetricLabel(location, "location"))

	for i := 0; i < 2; i++ {
		user := <-ch
		switch metrics.MetricLabel(user, "user") {
		case "foo":
			assert.Equal(t, 2.0, metrics.MetricValue(user).GetGauge().GetValue())
		case "bar":
			assert.Equal(t, 1.0, metrics.MetricValue(user).GetGauge().GetValue())
		default:
			t.Log("invalid user")
			t.Fail()
		}
	}
}

func TestCollector_Collect_Fail(t *testing.T) {
	c := plex.NewCollector("", "", "", time.Minute)
	c.(*plex.Collector).API = &client{failing: true}

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
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
