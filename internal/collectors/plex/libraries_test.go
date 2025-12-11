package plex

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLibraryCollector_Collect(t *testing.T) {
	tests := []struct {
		name   string
		getter libraryGetter
		want   string
	}{
		{
			name: "movie",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "movies", Type: "movie", Key: "1"}},
				movies: []plex.Movie{
					{Title: "movie 1", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
					{Title: "movie 2", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 2 * 1024}}}}},
				},
			},
			want: `
			# HELP mediamon_plex_library_bytes Library size in bytes
			# TYPE mediamon_plex_library_bytes gauge
			mediamon_plex_library_bytes{library="movies",url="http://localhost:8080"} 3072
			# HELP mediamon_plex_library_count Library size in number of entries
			# TYPE mediamon_plex_library_count gauge
			mediamon_plex_library_count{library="movies",url="http://localhost:8080"} 2
`,
		},
		{
			name: "movie - empty",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "movies", Type: "movie", Key: "1"}},
			},
			want: `
			# HELP mediamon_plex_library_bytes Library size in bytes
			# TYPE mediamon_plex_library_bytes gauge
			mediamon_plex_library_bytes{library="movies",url="http://localhost:8080"} 0
			# HELP mediamon_plex_library_count Library size in number of entries
			# TYPE mediamon_plex_library_count gauge
			mediamon_plex_library_count{library="movies",url="http://localhost:8080"} 0
			`,
		},
		{
			name: "show",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "shows", Type: "show", Key: "2"}},
				shows:     []plex.Show{{Key: "20", RatingKey: "21", Title: "show 1"}},
				seasons:   map[string][]plex.Season{"21": {{Key: "22", RatingKey: "23", Title: "Season 1"}}},
				episodes:  map[string][]plex.Episode{"23": {{Title: "Pilot", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}}}},
			},
			want: `
			# HELP mediamon_plex_library_bytes Library size in bytes
			# TYPE mediamon_plex_library_bytes gauge
			mediamon_plex_library_bytes{library="shows",url="http://localhost:8080"} 1024
			# HELP mediamon_plex_library_count Library size in number of entries
			# TYPE mediamon_plex_library_count gauge
			mediamon_plex_library_count{library="shows",url="http://localhost:8080"} 1
			`,
		},
		{
			name: "show - empty season",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "shows", Type: "show", Key: "2"}},
				shows:     []plex.Show{{Key: "20", RatingKey: "21", Title: "show 1"}},
				seasons:   map[string][]plex.Season{"21": {{Key: "22", RatingKey: "23", Title: "Season 1"}}},
			},
			want: `
			# HELP mediamon_plex_library_bytes Library size in bytes
			# TYPE mediamon_plex_library_bytes gauge
			mediamon_plex_library_bytes{library="shows",url="http://localhost:8080"} 0
			# HELP mediamon_plex_library_count Library size in number of entries
			# TYPE mediamon_plex_library_count gauge
			mediamon_plex_library_count{library="shows",url="http://localhost:8080"} 0
			`,
		},
		{
			name: "show - empty seasons",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "shows", Type: "show", Key: "2"}},
				shows:     []plex.Show{{Key: "20", RatingKey: "21", Title: "show 1"}},
			},
			want: `
			# HELP mediamon_plex_library_bytes Library size in bytes
			# TYPE mediamon_plex_library_bytes gauge
			mediamon_plex_library_bytes{library="shows",url="http://localhost:8080"} 0
			# HELP mediamon_plex_library_count Library size in number of entries
			# TYPE mediamon_plex_library_count gauge
			mediamon_plex_library_count{library="shows",url="http://localhost:8080"} 0
			`,
		},
		{
			name: "show - empty",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "shows", Type: "show", Key: "2"}},
			},
			want: `
			# HELP mediamon_plex_library_bytes Library size in bytes
			# TYPE mediamon_plex_library_bytes gauge
			mediamon_plex_library_bytes{library="shows",url="http://localhost:8080"} 0
			# HELP mediamon_plex_library_count Library size in number of entries
			# TYPE mediamon_plex_library_count gauge
			mediamon_plex_library_count{library="shows",url="http://localhost:8080"} 0
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newLibraryCollector(tt.getter, "http://localhost:8080", slog.New(slog.DiscardHandler))
			assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(tt.want)))
		})
	}
}
