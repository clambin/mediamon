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
				movies: []plex.MediaMetadata{
					{Title: "movie 1", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
					{Title: "movie 2", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 2 * 1024}}}}},
				},
			},
			want: `
# HELP mediamon_plex_movie_count Total number of movies in Plex library
# TYPE mediamon_plex_movie_count gauge
mediamon_plex_movie_count{library="movies",url="http://localhost:8080"} 2

# HELP mediamon_plex_library_bytes Library size in bytes
# TYPE mediamon_plex_library_bytes gauge
mediamon_plex_library_bytes{library="movies",url="http://localhost:8080"} 3072
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
			`,
		},
		{
			name: "show",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "shows", Type: "show", Key: "2"}},
				episodes: []plex.MediaMetadata{
					{Title: "Pilot", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
					{Title: "EP2", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
				},
			},
			want: `
# HELP mediamon_plex_show_count Total number of shows in Plex library
# TYPE mediamon_plex_show_count gauge
mediamon_plex_show_count{library="shows",url="http://localhost:8080"} 1

# HELP mediamon_plex_episode_count Total number of episodes in Plex library
# TYPE mediamon_plex_episode_count gauge
mediamon_plex_episode_count{library="shows",url="http://localhost:8080"} 2

# HELP mediamon_plex_library_bytes Library size in bytes
# TYPE mediamon_plex_library_bytes gauge
mediamon_plex_library_bytes{library="shows",url="http://localhost:8080"} 2048
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
