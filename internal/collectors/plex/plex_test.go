package plex_test

import (
	plexClient "github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := plex.NewCollector("1.0", "http://localhost:8888", "username", "password")
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
	c := plex.NewCollector("1.0", "", "", "")
	l := mocks.NewIPLocator(t)
	p := mocks.NewAPI(t)
	c.API = p
	c.IPLocator = l

	l.On("Locate", "1.2.3.4").Return(10.0, 20.0, nil)

	idResp := plexClient.Identity{}
	idResp.Version = "foo"
	p.On("GetIdentity", mock.AnythingOfType("*context.emptyCtx")).Return(idResp, nil)

	var sessions = plexClient.Sessions{}
	sessions.Metadata = []plexClient.Session{
		{
			Title:      "foo",
			Type:       "movie",
			Duration:   100,
			ViewOffset: 50,
			User:       plexClient.SessionUser{Title: "bar"},
			Player:     plexClient.SessionPlayer{Product: "Plex Web", Address: "192.168.0.1"},
			Media:      []plexClient.SessionMedia{{Part: []plexClient.MediaSessionPart{{Decision: "directplay"}}}},
			Session:    plexClient.SessionStats{ID: "1", Location: "lan"},
		},
		{
			GrandparentTitle: "foo",
			ParentIndex:      1,
			Index:            10,
			Title:            "bar",
			Type:             "episode",
			Duration:         100,
			ViewOffset:       75,
			User:             plexClient.SessionUser{Title: "bar"},
			Player:           plexClient.SessionPlayer{Product: "Plex Web", Address: "1.2.3.4"},
			Session:          plexClient.SessionStats{ID: "2", Location: "wan"},
			Media:            []plexClient.SessionMedia{{Part: []plexClient.MediaSessionPart{{Decision: "transcode"}}}},
			TranscodeSession: plexClient.SessionTranscoder{
				VideoDecision: "transcode",
				Speed:         21.0,
			},
		},
		{
			Title:      "foo",
			Type:       "movie",
			Duration:   100,
			ViewOffset: 10,
			User:       plexClient.SessionUser{Title: "bar"},
			Player:     plexClient.SessionPlayer{Product: "Plex Web", Address: "1.2.3.4"},
			Session:    plexClient.SessionStats{ID: "3", Location: "wan"},
			Media:      []plexClient.SessionMedia{{Part: []plexClient.MediaSessionPart{{Decision: "transcode"}}}},
			TranscodeSession: plexClient.SessionTranscoder{
				VideoDecision: "transcode",
				Throttled:     true,
			},
		},
	}

	p.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(sessions, nil)

	e := strings.NewReader(`# HELP mediamon_plex_session_count Active Plex session
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo - S01E10 - bar",url="",user="bar"} 0.75
mediamon_plex_session_count{address="1.2.3.4",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo",url="",user="bar"} 0.1
mediamon_plex_session_count{address="192.168.0.1",lat="",location="lan",lon="",mode="directplay",player="Plex Web",title="foo",url="",user="bar"} 0.5
# HELP mediamon_plex_transcoder_count Video transcode session
# TYPE mediamon_plex_transcoder_count gauge
mediamon_plex_transcoder_count{state="throttled",url=""} 1
mediamon_plex_transcoder_count{state="transcoding",url=""} 1
# HELP mediamon_plex_transcoder_speed Speed of active transcoder
# TYPE mediamon_plex_transcoder_speed gauge
mediamon_plex_transcoder_speed{url=""} 21
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="",version="foo"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(c, e))
}
