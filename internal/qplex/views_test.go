package qplex_test

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/qplex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetViews(t *testing.T) {
	type args struct {
		tokens  []string
		reverse bool
	}
	tests := []struct {
		name    string
		args    args
		want    []qplex.ViewCountEntry
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ascending",
			args: args{
				tokens:  []string{"1"},
				reverse: false,
			},
			want: []qplex.ViewCountEntry{
				{Library: "Movies", Title: "foo", Views: 1},
				{Library: "Shows", Title: "bar", Views: 2},
			},
			wantErr: assert.NoError,
		},
		{
			name: "descending",
			args: args{
				tokens:  []string{"1"},
				reverse: true,
			},
			want: []qplex.ViewCountEntry{
				{Library: "Shows", Title: "bar", Views: 2},
				{Library: "Movies", Title: "foo", Views: 1},
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiple tokens",
			args: args{
				tokens:  []string{"1", "2", "3"},
				reverse: true,
			},
			want: []qplex.ViewCountEntry{
				{Library: "Shows", Title: "bar", Views: 6},
				{Library: "Movies", Title: "foo", Views: 3},
			},
			wantErr: assert.NoError,
		},
	}
	c := fakeClient{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := qplex.GetViews(context.Background(), c, tt.args.tokens, tt.args.reverse)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type fakeClient struct {
}

var _ qplex.PlexGetter = &fakeClient{}

func (f fakeClient) SetAuthToken(_ string) {
}

func (f fakeClient) GetLibraries(_ context.Context) (plex.Libraries, error) {
	return plex.Libraries{
		Directory: []plex.LibrariesDirectory{
			{
				Key:   "1",
				Type:  "movie",
				Title: "Movies",
			},
			{
				Key:   "2",
				Type:  "show",
				Title: "Shows",
			},
		},
	}, nil
}

func (f fakeClient) GetMovieLibrary(_ context.Context, s string) (plex.MovieLibrary, error) {
	if s != "1" {
		return plex.MovieLibrary{}, fmt.Errorf("invalid movie library key: %s", s)
	}
	return plex.MovieLibrary{
		Metadata: []plex.MovieLibraryEntry{
			{
				Guid:      "1",
				Title:     "foo",
				ViewCount: 1,
			},
		},
	}, nil
}

func (f fakeClient) GetShowLibrary(_ context.Context, s string) (plex.ShowLibrary, error) {
	if s != "2" {
		return plex.ShowLibrary{}, fmt.Errorf("invalid show library key: %s", s)
	}
	return plex.ShowLibrary{
		Metadata: []plex.ShowLibraryEntry{
			{
				Guid:      "2",
				Title:     "bar",
				ViewCount: 2,
			},
		},
	}, nil
}
