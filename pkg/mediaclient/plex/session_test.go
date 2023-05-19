package plex_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
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

func TestSession_GetTitle(t *testing.T) {
	tests := []struct {
		name    string
		session plex.Session
		want    string
	}{
		{
			name: "movie",
			session: plex.Session{
				GrandparentTitle: "foo",
				ParentIndex:      1,
				Index:            1,
				Title:            "bar",
				Type:             "movie",
			},
			want: "bar",
		},
		{
			name: "episode",
			session: plex.Session{
				GrandparentTitle: "foo",
				ParentIndex:      1,
				Index:            10,
				Title:            "bar",
				Type:             "episode",
			},
			want: "foo - S01E10 - bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.session.GetTitle())
		})
	}
}

func TestSession_GetProgress(t *testing.T) {
	type fields struct {
		Duration   int
		ViewOffset int
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "start",
			fields: fields{
				Duration:   100,
				ViewOffset: 0,
			},
			want: 0,
		},
		{
			name: "half",
			fields: fields{
				Duration:   100,
				ViewOffset: 50,
			},
			want: 0.5,
		},
		{
			name: "full",
			fields: fields{
				Duration:   100,
				ViewOffset: 100,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := plex.Session{
				Duration:   tt.fields.Duration,
				ViewOffset: tt.fields.ViewOffset,
			}
			assert.Equal(t, tt.want, s.GetProgress())
		})
	}
}
