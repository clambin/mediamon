package plex

import (
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessSessions(t *testing.T) {
	var sessions = plex.SessionsResponse{}
	sessions.MediaContainer.Metadata = []plex.SessionsResponseRecord{
		{
			Title:   "foo",
			User:    plex.SessionsResponseRecordUser{Title: "bar"},
			Player:  plex.SessionsResponseRecordPlayer{Product: "Plex Web"},
			Session: plex.SessionsResponseRecordSession{ID: "1", Location: "lan"},
		},
		{
			Title:   "foo",
			User:    plex.SessionsResponseRecordUser{Title: "bar"},
			Player:  plex.SessionsResponseRecordPlayer{Product: "Plex Web"},
			Session: plex.SessionsResponseRecordSession{ID: "2", Location: "wan"},
			TranscodeSession: plex.SessionsResponseRecordTranscodeSession{
				VideoDecision: "transcode",
				Speed:         10.0,
			},
		},
		{
			Title:   "foo",
			User:    plex.SessionsResponseRecordUser{Title: "bar"},
			Player:  plex.SessionsResponseRecordPlayer{Product: "Plex Web"},
			Session: plex.SessionsResponseRecordSession{ID: "3", Location: "wan"},
			TranscodeSession: plex.SessionsResponseRecordTranscodeSession{
				VideoDecision: "transcode",
				Throttled:     true,
			},
		},
	}

	stats := parseSessions(sessions)
	assert.Len(t, stats, len(sessions.MediaContainer.Metadata))

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
