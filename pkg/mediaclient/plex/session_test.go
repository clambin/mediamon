package plex_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlexClient_GetStats(t *testing.T) {
	c, s := makeClientAndServer(nil)
	defer s.Close()

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

func TestSession_GetMediaMode(t *testing.T) {
	// media.part.decisions="[{parts:[{decision:directplay streamInfo:[{decision: kind:1 codec:h264 location:direct} {decision: kind:2 codec:eac3 location:direct} {decision:ignore kind:3 codec:srt location:}]}]}]" transcode.videoDecision=""
	// media.part.decisions="[{parts:[{decision:transcode streamInfo:[{decision:transcode kind:1 codec:h264 location:segments-av} {decision:transcode kind:2 codec:aac location:segments-av} {decision:transcode kind:3 codec:webvtt location:segments-subs}]}]}]" transcode.videoDecision=transcode
	// media.part.decisions="[{parts:[{decision:transcode streamInfo:[{decision:copy kind:1 codec:h264 location:segments-video} {decision:copy kind:2 codec:aac location:segments-audio} {decision:transcode kind:3 codec:ass location:segments-subs}]}]}]" transcode.videoDecision=copy
	// media.part.decisions="[{parts:[{decision: streamInfo:[{decision: kind:1 codec:h264 location:} {decision: kind:2 codec:aac location:} {decision: kind:3 codec:srt location:} {decision: kind:3 codec:srt location:}]}]}]" transcode.videoDecision=""
	type fields struct {
		media     []plex.SessionMedia
		transcode plex.SessionTranscoder
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "directplay",
			fields: fields{
				media: []plex.SessionMedia{{Part: []plex.MediaSessionPart{{Decision: "directplay"}}}},
			},
			want: "directplay",
		},
		{
			name: "copy",
			fields: fields{
				media:     []plex.SessionMedia{{Part: []plex.MediaSessionPart{{Decision: "transcode"}}}},
				transcode: plex.SessionTranscoder{VideoDecision: "copy"},
			},
			want: "copy",
		},
		{
			name: "transcode",
			fields: fields{
				media:     []plex.SessionMedia{{Part: []plex.MediaSessionPart{{Decision: "transcode"}}}},
				transcode: plex.SessionTranscoder{VideoDecision: "transcode"},
			},
			want: "transcode",
		},
		{
			name: "unknown",
			fields: fields{
				media: []plex.SessionMedia{{Part: []plex.MediaSessionPart{{}}}},
			},
			want: "unknown",
		},
		{
			name: "empty",
			fields: fields{
				media: []plex.SessionMedia{},
			},
			want: "unknown",
		},
		{
			name: "multiple",
			fields: fields{
				media: []plex.SessionMedia{
					{Part: []plex.MediaSessionPart{{Decision: "directplay"}}},
					{Part: []plex.MediaSessionPart{{Decision: "transcode"}}},
				},
				transcode: plex.SessionTranscoder{VideoDecision: "transcode"},
			},
			want: "directplay,transcode",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := plex.Session{
				Media:            tt.fields.media,
				TranscodeSession: tt.fields.transcode,
			}
			assert.Equal(t, tt.want, s.GetMediaMode())
		})
	}
}
