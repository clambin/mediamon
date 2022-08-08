package plex_test

import (
	"fmt"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/pkg/iplocator/mocks"
	plexAPI "github.com/clambin/mediamon/pkg/mediaclient/plex"
	plexMock "github.com/clambin/mediamon/pkg/mediaclient/plex/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
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

	e := strings.NewReader(`# HELP mediamon_plex_session_count Active Plex session
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",id="2",lat="20.00",location="wan",lon="10.00",player="Plex Web",title="foo",url="",user="bar"} 1
mediamon_plex_session_count{address="1.2.3.4",id="3",lat="20.00",location="wan",lon="10.00",player="Plex Web",title="foo",url="",user="bar"} 1
mediamon_plex_session_count{address="192.168.0.1",id="1",lat="",location="lan",lon="",player="Plex Web",title="foo",url="",user="bar"} 1
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

func TestCollector_Collect_Fail(t *testing.T) {
	c := plex.NewCollector("", "", "")
	m := &plexMock.API{}
	c.(*plex.Collector).API = m

	m.On("GetIdentity", mock.AnythingOfType("*context.emptyCtx")).Return(plexAPI.IdentityResponse{}, fmt.Errorf("failure"))
	m.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(plexAPI.SessionsResponse{}, fmt.Errorf("failure"))

	err := testutil.CollectAndCompare(c, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `Desc{fqName: "mediamon_error", help: "Error getting Plex version", constLabels: {}, variableLabels: []}`)
}
