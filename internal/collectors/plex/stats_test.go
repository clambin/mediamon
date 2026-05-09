package plex

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestStatsCollector_Collect(t *testing.T) {
	tests := []struct {
		name   string
		getter statsGetter
		want   string
	}{
		{
			name:   "empty library",
			getter: fakeGetter{},
			want: `
# HELP mediamon_plex_episode_count Total number of episodes in Plex library
# TYPE mediamon_plex_episode_count gauge
mediamon_plex_episode_count{url="http://localhost"} 0

# HELP mediamon_plex_movie_count Total number of movies in Plex library
# TYPE mediamon_plex_movie_count gauge
mediamon_plex_movie_count{url="http://localhost"} 0

# HELP mediamon_plex_show_count Total number of shows in Plex library
# TYPE mediamon_plex_show_count gauge
mediamon_plex_show_count{url="http://localhost"} 0
`,
		},
		{
			name: "movies",
			getter: fakeGetter{
				libraries: []plex.Library{{Title: "movies", Type: "movie", Key: "1"}},
				movies: []plex.Movie{
					{Title: "movie 1", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
					{Title: "movie 2", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 2 * 1024}}}}},
				},
			},
			want: `
# HELP mediamon_plex_episode_count Total number of episodes in Plex library
# TYPE mediamon_plex_episode_count gauge
mediamon_plex_episode_count{url="http://localhost"} 0

# HELP mediamon_plex_movie_count Total number of movies in Plex library
# TYPE mediamon_plex_movie_count gauge
mediamon_plex_movie_count{url="http://localhost"} 2

# HELP mediamon_plex_show_count Total number of shows in Plex library
# TYPE mediamon_plex_show_count gauge
mediamon_plex_show_count{url="http://localhost"} 0
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
# HELP mediamon_plex_episode_count Total number of episodes in Plex library
# TYPE mediamon_plex_episode_count gauge
mediamon_plex_episode_count{url="http://localhost"} 1

# HELP mediamon_plex_movie_count Total number of movies in Plex library
# TYPE mediamon_plex_movie_count gauge
mediamon_plex_movie_count{url="http://localhost"} 0

# HELP mediamon_plex_show_count Total number of shows in Plex library
# TYPE mediamon_plex_show_count gauge
mediamon_plex_show_count{url="http://localhost"} 1
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newStatsCollector(tt.getter, "http://localhost", slog.New(slog.DiscardHandler))
			assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(tt.want)))
		})
	}
}
