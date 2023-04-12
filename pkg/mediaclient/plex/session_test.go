package plex_test

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexClient_GetStats(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()

	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	c := plex.Client{
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	sessions, err := c.GetSessions(context.Background())
	require.NoError(t, err)

	titles := []string{"pilot", "movie 1", "movie 2", "movie 3"}
	locations := []string{"lan", "wan", "lan", "lan"}
	require.Len(t, sessions.Metadata, len(titles))

	for index, title := range titles {
		assert.Equal(t, title, sessions.Metadata[index].Title)
		assert.Equal(t, "Plex Web", sessions.Metadata[index].Player.Product)
		assert.Equal(t, locations[index], sessions.Metadata[index].Session.Location)

		if sessions.Metadata[index].TranscodeSession.VideoDecision == "transcode" {
			assert.NotZero(t, sessions.Metadata[index].TranscodeSession.Speed)
		} else {
			assert.Zero(t, sessions.Metadata[index].TranscodeSession.Speed)
		}
	}
}
