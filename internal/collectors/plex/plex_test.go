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

func TestCollector_Collect(t *testing.T) {
	c := plex.NewCollector("1.0", "", "", "")
	l := mocks.NewIPLocator(t)
	p := mocks.NewGetter(t)
	c.Plex = p
	c.IPLocator = l

	l.EXPECT().Locate("1.2.3.4").Return(10.0, 20.0, nil)
	p.EXPECT().GetIdentity(mock.Anything).Return(plexClient.Identity{Version: "foo"}, nil)

	testCases := []struct {
		name    string
		session plexClient.Session
		want    string
	}{
		{
			name: "direct",
			session: plexClient.Session{
				Title:      "foo",
				Type:       "movie",
				Duration:   100,
				ViewOffset: 50,
				User:       plexClient.SessionUser{Title: "bar"},
				Player:     plexClient.SessionPlayer{Product: "Plex Web", Address: "192.168.0.1"},
				Media:      []plexClient.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac", Part: []plexClient.MediaSessionPart{{Decision: "directplay"}}}},
				Session:    plexClient.SessionStats{ID: "1", Location: "lan"},
			},
			want: `
# HELP mediamon_plex_session_count Active Plex session
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="192.168.0.1",audioCodec="aac",lat="",location="lan",lon="",mode="directplay",player="Plex Web",title="foo",url="",user="bar",videoCodec="hvec"} 0.5
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="",version="foo"} 1
`,
		},
		{
			name: "transcode",
			session: plexClient.Session{
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
				Media:            []plexClient.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac", Part: []plexClient.MediaSessionPart{{Decision: "transcode"}}}},
				TranscodeSession: plexClient.SessionTranscoder{
					VideoDecision: "transcode",
					Speed:         21.0,
				},
			},
			want: `
# HELP mediamon_plex_session_count Active Plex session
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",audioCodec="aac",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo - S01E10 - bar",url="",user="bar",videoCodec="hvec"} 0.75
# HELP mediamon_plex_transcoder_count Video transcode session
# TYPE mediamon_plex_transcoder_count gauge
mediamon_plex_transcoder_count{state="throttled",url=""} 0
mediamon_plex_transcoder_count{state="transcoding",url=""} 1
# HELP mediamon_plex_transcoder_speed Speed of active transcoder
# TYPE mediamon_plex_transcoder_speed gauge
mediamon_plex_transcoder_speed{url=""} 21
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="",version="foo"} 1
`,
		},
		{
			name: "throttled",
			session: plexClient.Session{
				Title:      "foo",
				Type:       "movie",
				Duration:   100,
				ViewOffset: 10,
				User:       plexClient.SessionUser{Title: "bar"},
				Player:     plexClient.SessionPlayer{Product: "Plex Web", Address: "1.2.3.4"},
				Session:    plexClient.SessionStats{ID: "3", Location: "wan"},
				Media:      []plexClient.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac", Part: []plexClient.MediaSessionPart{{Decision: "transcode"}}}},
				TranscodeSession: plexClient.SessionTranscoder{
					VideoDecision: "transcode",
					Throttled:     true,
				},
			},
			want: `
# HELP mediamon_plex_session_count Active Plex session
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",audioCodec="aac",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo",url="",user="bar",videoCodec="hvec"} 0.1
# HELP mediamon_plex_transcoder_count Video transcode session
# TYPE mediamon_plex_transcoder_count gauge
mediamon_plex_transcoder_count{state="throttled",url=""} 1
mediamon_plex_transcoder_count{state="transcoding",url=""} 0
# HELP mediamon_plex_transcoder_speed Speed of active transcoder
# TYPE mediamon_plex_transcoder_speed gauge
mediamon_plex_transcoder_speed{url=""} 0
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="",version="foo"} 1
`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p.EXPECT().GetSessions(mock.Anything).Return(plexClient.Sessions{Size: 1, Metadata: []plexClient.Session{tt.session}}, nil).Once()
			e := strings.NewReader(tt.want)

			r := prometheus.NewPedanticRegistry()
			r.MustRegister(c)

			assert.NoError(t, testutil.GatherAndCompare(r, e))
		})
	}
}
