package plex

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	p := fakeGetter{
		libraries: []plex.Library{
			{Title: "movies", Type: "movie", Key: "1"},
			{Title: "shows", Type: "show", Key: "2"},
		},
		movies: []plex.Movie{
			{Title: "a movie", Key: "10", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
		},
		identity: plex.Identity{Version: "1.0"},
	}
	c := NewCollector("1.0", "http://localhost:8080", "", "", http.DefaultClient, slog.New(slog.DiscardHandler))
	c.libraryCollector.(*libraryCollector).libraryGetter = p
	c.versionCollector.identityGetter = p
	c.sessionCollector.sessionGetter = p

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
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
`), "mediamon_plex_library_bytes", "mediamon_plex_library_count", "mediamon_plex_version"))
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
