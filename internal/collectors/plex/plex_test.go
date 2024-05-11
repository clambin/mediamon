package plex

import (
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	p := mocks.NewGetter(t)
	p.EXPECT().GetIdentity(mock.Anything).Return(plex.Identity{Version: "1.0"}, nil)
	p.EXPECT().GetSessions(mock.Anything).Return([]plex.Session{}, nil)
	p.EXPECT().GetLibraries(mock.Anything).Return([]plex.Library{
		{Title: "movies", Type: "movie", Key: "1"},
		{Title: "shows", Type: "show", Key: "2"},
	}, nil)
	p.EXPECT().GetMovies(mock.Anything, "1").Return([]plex.Movie{
		{Title: "a movie", Key: "10", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
	}, nil)
	p.EXPECT().GetShows(mock.Anything, "2").Return([]plex.Show{}, nil)

	cb := NewCollector("1.0", "http://localhost:8080", "", "", slog.Default())
	cb.Collector.(*Collector).libraryCollector.libraryGetter = p
	cb.Collector.(*Collector).versionCollector.versionGetter = p
	cb.Collector.(*Collector).sessionCollector.sessionGetter = p

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(cb)

	assert.NoError(t, testutil.GatherAndCompare(r, strings.NewReader(`
# HELP mediamon_plex_library_bytes Library size in bytes
# TYPE mediamon_plex_library_bytes gauge
mediamon_plex_library_bytes{library="movies",url="http://localhost:8080"} 1024
mediamon_plex_library_bytes{library="shows",url="http://localhost:8080"} 0
# HELP mediamon_plex_library_count Library size in number of entries
# TYPE mediamon_plex_library_count gauge
mediamon_plex_library_count{library="movies",url="http://localhost:8080"} 1
mediamon_plex_library_count{library="shows",url="http://localhost:8080"} 0
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="http://localhost:8080",version="1.0"} 1
`)))
}

func Test_chopPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		///library/metadata", "/library/sections
		{
			name: "unaltered",
			path: "/healthz",
			want: "/healthz",
		},
		{
			name: "metadata",
			path: "/library/metadata/123/details",
			want: "/library/metadata",
		},
		{
			name: "sections",
			path: "/library/sections/123/details",
			want: "/library/sections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out := chopPath(&http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: tt.path},
			})
			assert.Equal(t, http.MethodPost, out.Method)
			assert.Equal(t, tt.want, out.URL.Path)
		})
	}
}
