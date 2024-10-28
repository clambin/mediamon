package plex

import (
	"github.com/clambin/mediaclients/plex"
	collectorbreaker "github.com/clambin/mediamon/v2/collector-breaker"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/clambin/mediamon/v2/iplocator"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"strings"
	"testing"
)

func TestSessionsCollector_Collector(t *testing.T) {
	tests := []struct {
		name    string
		session plex.Session
		want    string
	}{
		{
			name: "direct",
			session: plex.Session{
				Title:      "foo",
				Type:       "movie",
				Duration:   100,
				ViewOffset: 50,
				User:       plex.SessionUser{Title: "bar"},
				Player:     plex.SessionPlayer{Product: "Plex Web", Address: "192.168.0.1"},
				Media:      []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac", Part: []plex.MediaSessionPart{{Decision: "directplay"}}}},
				Session:    plex.SessionStats{ID: "1", Location: "lan"},
			},
			want: `
# HELP mediamon_plex_session_bandwidth Active Plex session Bandwidth usage (in kbps)
# TYPE mediamon_plex_session_bandwidth gauge
mediamon_plex_session_bandwidth{address="192.168.0.1",audioCodec="aac",lat="",location="lan",lon="",mode="directplay",player="Plex Web",title="foo",url="http://localhost:8080",user="bar",videoCodec="hvec"} 0
# HELP mediamon_plex_session_count Active Plex session progress
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="192.168.0.1",audioCodec="aac",lat="",location="lan",lon="",mode="directplay",player="Plex Web",title="foo",url="http://localhost:8080",user="bar",videoCodec="hvec"} 0.5
`,
		},
		{
			name: "transcode",
			session: plex.Session{
				GrandparentTitle: "foo",
				ParentIndex:      1,
				Index:            10,
				Title:            "bar",
				Type:             "episode",
				Duration:         100,
				ViewOffset:       75,
				User:             plex.SessionUser{Title: "bar"},
				Player:           plex.SessionPlayer{Product: "Plex Web", Address: "1.2.3.4"},
				Session:          plex.SessionStats{ID: "2", Location: "wan"},
				Media:            []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac", Part: []plex.MediaSessionPart{{Decision: "transcode"}}}},
				TranscodeSession: plex.SessionTranscoder{VideoDecision: "transcode", Speed: 21.0},
			},
			want: `
# HELP mediamon_plex_session_bandwidth Active Plex session Bandwidth usage (in kbps)
# TYPE mediamon_plex_session_bandwidth gauge
mediamon_plex_session_bandwidth{address="1.2.3.4",audioCodec="aac",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo - S01E10 - bar",url="http://localhost:8080",user="bar",videoCodec="hvec"} 0
# HELP mediamon_plex_session_count Active Plex session progress
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",audioCodec="aac",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo - S01E10 - bar",url="http://localhost:8080",user="bar",videoCodec="hvec"} 0.75
# HELP mediamon_plex_transcoder_count Video transcode session
# TYPE mediamon_plex_transcoder_count gauge
mediamon_plex_transcoder_count{state="throttled",url="http://localhost:8080"} 0
mediamon_plex_transcoder_count{state="transcoding",url="http://localhost:8080"} 1
# HELP mediamon_plex_transcoder_speed Speed of active transcoder
# TYPE mediamon_plex_transcoder_speed gauge
mediamon_plex_transcoder_speed{url="http://localhost:8080"} 21
`,
		},
		{
			name: "transcode - throttled",
			session: plex.Session{
				GrandparentTitle: "foo",
				ParentIndex:      1,
				Index:            10,
				Title:            "bar",
				Type:             "episode",
				Duration:         100,
				ViewOffset:       75,
				User:             plex.SessionUser{Title: "bar"},
				Player:           plex.SessionPlayer{Product: "Plex Web", Address: "1.2.3.4"},
				Session:          plex.SessionStats{ID: "2", Location: "wan"},
				Media:            []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac", Part: []plex.MediaSessionPart{{Decision: "transcode"}}}},
				TranscodeSession: plex.SessionTranscoder{VideoDecision: "transcode", Speed: 21.0, Throttled: true},
			},
			want: `
# HELP mediamon_plex_session_bandwidth Active Plex session Bandwidth usage (in kbps)
# TYPE mediamon_plex_session_bandwidth gauge
mediamon_plex_session_bandwidth{address="1.2.3.4",audioCodec="aac",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo - S01E10 - bar",url="http://localhost:8080",user="bar",videoCodec="hvec"} 0
# HELP mediamon_plex_session_count Active Plex session progress
# TYPE mediamon_plex_session_count gauge
mediamon_plex_session_count{address="1.2.3.4",audioCodec="aac",lat="20.00",location="wan",lon="10.00",mode="transcode",player="Plex Web",title="foo - S01E10 - bar",url="http://localhost:8080",user="bar",videoCodec="hvec"} 0.75
# HELP mediamon_plex_transcoder_count Video transcode session
# TYPE mediamon_plex_transcoder_count gauge
mediamon_plex_transcoder_count{state="throttled",url="http://localhost:8080"} 1
mediamon_plex_transcoder_count{state="transcoding",url="http://localhost:8080"} 0
# HELP mediamon_plex_transcoder_speed Speed of active transcoder
# TYPE mediamon_plex_transcoder_speed gauge
mediamon_plex_transcoder_speed{url="http://localhost:8080"} 21
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := mocks.NewGetter(t)
			p.EXPECT().GetSessions(mock.Anything).Return([]plex.Session{tt.session}, nil).Once()
			i := mocks.NewIPLocator(t)
			if tt.session.Session.Location == "wan" {
				i.EXPECT().Locate("1.2.3.4").Return(iplocator.Location{Lon: 10, Lat: 20}, nil).Once()
			}

			c := sessionCollector{
				sessionGetter: p,
				IPLocator:     i,
				url:           "http://localhost:8080",
				logger:        slog.Default(),
			}
			assert.NoError(t, testutil.CollectAndCompare(
				collectorbreaker.PassThroughCollector{Collector: c},
				strings.NewReader(tt.want),
			))
		})
	}
}

/*
func Test_plexSessions(t *testing.T) {
	testCases := []struct {
		name    string
		session plex.Session
		want    map[string]plexSession
	}{
		{
			name: "direct",
			session: plex.Session{
				Title:   "foo",
				User:    plex.SessionUser{Title: "bar"},
				Player:  plex.SessionPlayer{Product: "Plex Web"},
				Media:   []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac"}},
				Session: plex.SessionStats{ID: "1", Location: "lan"},
			},
			want: map[string]plexSession{
				"1": {user: "bar", player: "Plex Web", location: "lan", title: "foo", videoMode: "unknown", audioCodec: "aac", videoCodec: "hvec"},
			},
		},
		{
			name: "transcode",
			session: plex.Session{
				Title:            "foo",
				User:             plex.SessionUser{Title: "bar"},
				Player:           plex.SessionPlayer{Product: "Plex Web"},
				Media:            []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac"}},
				Session:          plex.SessionStats{ID: "2", Location: "lan"},
				TranscodeSession: plex.SessionTranscoder{VideoDecision: "transcode", Speed: 10.0},
			},
			want: map[string]plexSession{
				"2": {user: "bar", player: "Plex Web", location: "lan", title: "foo", videoMode: "unknown", speed: 10, audioCodec: "aac", videoCodec: "hvec"},
			},
		},
		{
			name: "throttled",
			session: plex.Session{
				Title:            "foo",
				User:             plex.SessionUser{Title: "bar"},
				Player:           plex.SessionPlayer{Product: "Plex Web"},
				Media:            []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac"}},
				Session:          plex.SessionStats{ID: "3", Location: "lan"},
				TranscodeSession: plex.SessionTranscoder{VideoDecision: "transcode", Throttled: true},
			},
			want: map[string]plexSession{
				"3": {user: "bar", player: "Plex Web", location: "lan", title: "foo", videoMode: "unknown", throttled: true, audioCodec: "aac", videoCodec: "hvec"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var c Collector
			s := make(map[string]plexSession)
			for id, sess := range c.plexSessions([]plex.Session{tt.session}) {
				s[id] = sess
			}
			assert.Equal(t, tt.want, s)
		})
	}
}


*/
