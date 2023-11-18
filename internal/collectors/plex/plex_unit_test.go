package plex

import (
	"github.com/clambin/mediaclients/plex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessSessions(t *testing.T) {
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
				Session:          plex.SessionStats{ID: "2", Location: "wan"},
				TranscodeSession: plex.SessionTranscoder{VideoDecision: "transcode", Speed: 10.0},
			},
			want: map[string]plexSession{
				"2": {user: "bar", player: "Plex Web", location: "wan", title: "foo", videoMode: "unknown", speed: 10, audioCodec: "aac", videoCodec: "hvec"},
			},
		},
		{
			name: "throttled",
			session: plex.Session{
				Title:            "foo",
				User:             plex.SessionUser{Title: "bar"},
				Player:           plex.SessionPlayer{Product: "Plex Web"},
				Media:            []plex.SessionMedia{{VideoCodec: "hvec", AudioCodec: "aac"}},
				Session:          plex.SessionStats{ID: "3", Location: "wan"},
				TranscodeSession: plex.SessionTranscoder{VideoDecision: "transcode", Throttled: true},
			},
			want: map[string]plexSession{
				"3": {user: "bar", player: "Plex Web", location: "wan", title: "foo", videoMode: "unknown", throttled: true, audioCodec: "aac", videoCodec: "hvec"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			stats := parseSessions([]plex.Session{tt.session})
			assert.Equal(t, tt.want, stats)
		})
	}
}
