package plex

import (
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	g := fakeGetter{
		libraries: []plex.Library{
			{Title: "movies", Type: "movie", Key: "1"},
			{Title: "shows", Type: "show", Key: "2"},
		},
		movies: []plex.Movie{
			{Title: "a movie", Key: "10", Media: []plex.Media{{Part: []plex.MediaPart{{Size: 1024}}}}},
		},
		identity: plex.Identity{Version: "1.0"},
	}

	c := NewCollector(
		"http://localhost:8080",
		Config{Version: "1.0"},
		http.DefaultClient,
		slog.New(slog.DiscardHandler),
	)
	for _, coll := range c.collectors {
		switch cl := coll.(type) {
		case *libraryCollector:
			cl.libraryGetter = g
		case *versionCollector:
			cl.identityGetter = g
		case *sessionCollector:
			cl.sessionGetter = g
		default:
			panic("unknown collector")
		}
	}

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
