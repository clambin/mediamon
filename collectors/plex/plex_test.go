package plex_test

import (
	"fmt"
	"github.com/clambin/go-metrics/tools"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/pkg/iplocator/mocks"
	plexAPI "github.com/clambin/mediamon/pkg/mediaclient/plex"
	plexMock "github.com/clambin/mediamon/pkg/mediaclient/plex/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := plex.NewCollector("http://localhost:8888", "username", "password")
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{
		"mediamon_plex_version",
		"mediamon_plex_session_count",
		"mediamon_plex_transcoder_count",
		"mediamon_plex_transcoder_speed",
	} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"", metricName)
	}
}

func TestCollector_Collect(t *testing.T) {
	c := plex.NewCollector("", "", "")
	l := &mocks.Locator{}
	m := &plexMock.API{}
	c.(*plex.Collector).API = m
	c.(*plex.Collector).Locator = l

	l.On("Locate", "1.2.3.4").Return(10.0, 20.0, nil)

	idResp := plexAPI.IdentityResponse{}
	idResp.MediaContainer.Version = "foo"
	m.On("GetIdentity", mock.AnythingOfType("*context.emptyCtx")).Return(idResp, nil)

	var sessions = plexAPI.SessionsResponse{}
	sessions.MediaContainer.Metadata = []plexAPI.SessionsResponseRecord{
		{
			Title:   "foo",
			User:    plexAPI.SessionsResponseRecordUser{Title: "bar"},
			Player:  plexAPI.SessionsResponseRecordPlayer{Product: "Plex Web", Address: "192.168.0.1"},
			Session: plexAPI.SessionsResponseRecordSession{ID: "1", Location: "lan"},
		},
		{
			Title:   "foo",
			User:    plexAPI.SessionsResponseRecordUser{Title: "bar"},
			Player:  plexAPI.SessionsResponseRecordPlayer{Product: "Plex Web", Address: "1.2.3.4"},
			Session: plexAPI.SessionsResponseRecordSession{ID: "2", Location: "wan"},
			TranscodeSession: plexAPI.SessionsResponseRecordTranscodeSession{
				VideoDecision: "transcode",
				Speed:         21.0,
			},
		},
		{
			Title:   "foo",
			User:    plexAPI.SessionsResponseRecordUser{Title: "bar"},
			Player:  plexAPI.SessionsResponseRecordPlayer{Product: "Plex Web", Address: "1.2.3.4"},
			Session: plexAPI.SessionsResponseRecordSession{ID: "3", Location: "wan"},
			TranscodeSession: plexAPI.SessionsResponseRecordTranscodeSession{
				VideoDecision: "transcode",
				Throttled:     true,
			},
		},
	}

	m.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(sessions, nil)

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	version := <-ch
	assert.Equal(t, 1.0, tools.MetricValue(version).GetGauge().GetValue())
	assert.Equal(t, "foo", tools.MetricLabel(version, "version"))

	for i := 0; i < 3; i++ {
		session := <-ch

		require.Contains(t, session.Desc().String(), "mediamon_plex_session_count", i)
		assert.Equal(t, "bar", tools.MetricLabel(session, "user"), i)
		assert.Equal(t, "Plex Web", tools.MetricLabel(session, "player"), i)
		assert.Equal(t, "foo", tools.MetricLabel(session, "title"), i)
		assert.Equal(t, 1.0, tools.MetricValue(session).GetGauge().GetValue())

		var location string
		if tools.MetricLabel(session, "id") == "1" {
			location = "lan"
		} else {
			location = "wan"
		}
		assert.Equal(t, location, tools.MetricLabel(session, "location"))
		if location == "lan" {
			assert.Equal(t, "192.168.0.1", tools.MetricLabel(session, "address"))
			assert.Empty(t, tools.MetricLabel(session, "lon"))
			assert.Empty(t, tools.MetricLabel(session, "lat"))
		} else {
			assert.Equal(t, "1.2.3.4", tools.MetricLabel(session, "address"))
			assert.Equal(t, "10.00", tools.MetricLabel(session, "lon"))
			assert.Equal(t, "20.00", tools.MetricLabel(session, "lat"))
		}
	}

	for i := 0; i < 2; i++ {
		transcoder := <-ch
		require.Contains(t, transcoder.Desc().String(), "mediamon_plex_transcoder_count", i)
		assert.Equal(t, []string{"transcoding", "throttled"}[i], tools.MetricLabel(transcoder, "state"))
		assert.Equal(t, 1.0, tools.MetricValue(transcoder).GetGauge().GetValue())
	}

	speed := <-ch
	require.Contains(t, speed.Desc().String(), "mediamon_plex_transcoder_speed")
	assert.Equal(t, 21.0, tools.MetricValue(speed).GetGauge().GetValue())
}

func TestCollector_Collect_Fail(t *testing.T) {
	c := plex.NewCollector("", "", "")
	m := &plexMock.API{}
	c.(*plex.Collector).API = m

	m.On("GetIdentity", mock.AnythingOfType("*context.emptyCtx")).Return(plexAPI.IdentityResponse{}, fmt.Errorf("failure"))
	m.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(plexAPI.SessionsResponse{}, fmt.Errorf("failure"))
	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
	assert.Equal(t, `Desc{fqName: "mediamon_error", help: "Error getting Plex version", constLabels: {}, variableLabels: []}`, metric.Desc().String())
}
