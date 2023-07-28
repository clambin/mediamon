package plex

import (
	"github.com/clambin/mediaclients/plex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessSessions(t *testing.T) {
	var sessions = plex.Sessions{}
	sessions.Metadata = []plex.Session{
		{
			Title:   "foo",
			User:    plex.SessionUser{Title: "bar"},
			Player:  plex.SessionPlayer{Product: "Plex Web"},
			Session: plex.SessionStats{ID: "1", Location: "lan"},
		},
		{
			Title:   "foo",
			User:    plex.SessionUser{Title: "bar"},
			Player:  plex.SessionPlayer{Product: "Plex Web"},
			Session: plex.SessionStats{ID: "2", Location: "wan"},
			TranscodeSession: plex.SessionTranscoder{
				VideoDecision: "transcode",
				Speed:         10.0,
			},
		},
		{
			Title:   "foo",
			User:    plex.SessionUser{Title: "bar"},
			Player:  plex.SessionPlayer{Product: "Plex Web"},
			Session: plex.SessionStats{ID: "3", Location: "wan"},
			TranscodeSession: plex.SessionTranscoder{
				VideoDecision: "transcode",
				Throttled:     true,
			},
		},
	}

	stats := parseSessions(sessions)
	assert.Len(t, stats, len(sessions.Metadata))

	entry, ok := stats["1"]
	assert.True(t, ok)
	assert.Equal(t, "lan", entry.location)

	entry, ok = stats["2"]
	assert.True(t, ok)
	assert.Equal(t, "wan", entry.location)

	entry, ok = stats["3"]
	assert.True(t, ok)
	assert.True(t, entry.throttled)

}
