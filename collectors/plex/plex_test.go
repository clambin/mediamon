package plex_test

import (
	"github.com/clambin/mediamon/v2/collectors/plex"
	"github.com/clambin/mediamon/v2/collectors/plex/mocks"
	plexClient "github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
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
			Session:    plexClient.SessionStats{ID: "1", Location: "lan"},
		},
		{
			GrandparentTitle: "foo",
			ParentTitle:      "season 1",
			Title:            "bar",
			Type:             "episode",
			Duration:         100,
			ViewOffset:       75,
			User:             plexClient.SessionUser{Title: "bar"},
			Player:           plexClient.SessionPlayer{Product: "Plex Web", Address: "1.2.3.4"},
			Session:          plexClient.SessionStats{ID: "2", Location: "wan"},
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
			TranscodeSession: plexClient.SessionTranscoder{
				VideoDecision: "transcode",
				Throttled:     true,
			},
		},
	}

	p.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(sessions, nil)

	e := strings.NewReader(`# HELP mediamon_plex_session_count Active Plex session
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",id="2",lat="20.00",location="wan",lon="10.00",player="Plex Web",title="foo / season 1 / bar",url="",user="bar"} 0.75
mediamon_plex_session_count{address="1.2.3.4",id="3",lat="20.00",location="wan",lon="10.00",player="Plex Web",title="foo",url="",user="bar"} 0.1
mediamon_plex_session_count{address="192.168.0.1",id="1",lat="",location="lan",lon="",player="Plex Web",title="foo",url="",user="bar"} 0.5
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

/*
func TestCollector_Collect_Fail(t *testing.T) {
	c := plex.NewCollector("", "", "")
	m := &plexMock.API{}
	c.API = m

	m.On("GetIdentity", mock.AnythingOfType("*context.emptyCtx")).Return(plexClient.Identity{}, fmt.Errorf("failure"))
	m.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(plexClient.Sessions{}, fmt.Errorf("failure"))

	err := testutil.CollectAndCompare(c, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `Desc{fqName: "mediamon_error", help: "Error getting Plex version", constLabels: {}, variableLabels: []}`)
}
*/
