package plex

import (
	"errors"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"strings"
	"testing"
)

func TestLibraryCollector_Collect(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(p *mocks.Getter)
		want  string
	}{
		{
			name: "movie",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{{Title: "movies", Type: "movie", Key: "1"}}, nil)
				p.EXPECT().GetMovies(mock.Anything, "1").Return([]plex.Movie{
					{Title: "movie 1", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
					{Title: "movie 2", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 2 * 1024}}}}},
				}, nil)
			},
			want: `
# HELP mediamon_plex_library_entry_bytes Library file sizes
# TYPE mediamon_plex_library_entry_bytes gauge
mediamon_plex_library_entry_bytes{library="movies",title="movie 1",url="http://localhost:8080"} 1024
mediamon_plex_library_entry_bytes{library="movies",title="movie 2",url="http://localhost:8080"} 2048
`,
		},
		{
			name: "movie - empty",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{{Title: "movies", Type: "movie", Key: "1"}}, nil)
				p.EXPECT().GetMovies(mock.Anything, "1").Return([]plex.Movie{}, nil)
			},
			want: ``,
		},
		{
			name: "movie - error",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return(nil, errors.New("plex is down"))
			},
			want: ``,
		},
		{
			name: "show",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{{Title: "shows", Type: "show", Key: "2"}}, nil)
				p.EXPECT().GetShows(mock.Anything, "2").Return([]plex.Show{
					{Key: "20", RatingKey: "21", Title: "show 1"},
				}, nil)
				p.EXPECT().GetSeasons(mock.Anything, "21").Return([]plex.Season{
					{Key: "22", RatingKey: "23", Title: "Season 1"},
				}, nil)
				p.EXPECT().GetEpisodes(mock.Anything, "23").Return([]plex.Episode{
					{Title: "Pilot", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
				}, nil)
			},
			want: `
# HELP mediamon_plex_library_entry_bytes Library file sizes
# TYPE mediamon_plex_library_entry_bytes gauge
mediamon_plex_library_entry_bytes{library="shows",title="show 1",url="http://localhost:8080"} 1024
`,
		},
		{
			name: "show - empty season",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{{Title: "shows", Type: "show", Key: "2"}}, nil)
				p.EXPECT().GetShows(mock.Anything, "2").Return([]plex.Show{
					{Key: "20", RatingKey: "21", Title: "show 1"},
				}, nil)
				p.EXPECT().GetSeasons(mock.Anything, "21").Return([]plex.Season{
					{Key: "22", RatingKey: "23", Title: "Season 1"},
				}, nil)
				p.EXPECT().GetEpisodes(mock.Anything, "23").Return([]plex.Episode{}, nil)
			},
			want: ``,
		},
		{
			name: "show - empty seasons",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{{Title: "shows", Type: "show", Key: "2"}}, nil)
				p.EXPECT().GetShows(mock.Anything, "2").Return([]plex.Show{
					{Key: "20", RatingKey: "21", Title: "show 1"},
				}, nil)
				p.EXPECT().GetSeasons(mock.Anything, "21").Return([]plex.Season{}, nil)
			},
			want: ``,
		},
		{
			name: "show - empty",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{{Title: "shows", Type: "show", Key: "2"}}, nil)
				p.EXPECT().GetShows(mock.Anything, "2").Return([]plex.Show{}, nil)
			},
			want: ``,
		},
		{
			name: "show - error",
			setup: func(p *mocks.Getter) {
				p.EXPECT().GetLibraries(mock.Anything).Return(nil, errors.New("plex is down"))
			},
			want: ``,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := mocks.NewGetter(t)
			tt.setup(p)

			c := libraryCollector{
				libraryGetter: p,
				url:           "http://localhost:8080",
				l:             slog.Default(),
			}
			r := prometheus.NewPedanticRegistry()
			r.MustRegister(c)
			assert.NoError(t, testutil.GatherAndCompare(r, strings.NewReader(tt.want)))
		})
	}
}
