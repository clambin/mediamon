package plex

import (
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	c := NewCollector("1.0", "http://localhost:8080", "", "")
	c.libraryCollector.libraryGetter = p
	c.versionCollector.versionGetter = p
	c.sessionCollector.sessionGetter = p

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(c)

	assert.NoError(t, testutil.GatherAndCompare(r, strings.NewReader(`
# HELP mediamon_plex_library_entry_bytes Library file sizes
# TYPE mediamon_plex_library_entry_bytes gauge
mediamon_plex_library_entry_bytes{library="movies",title="a movie",url="http://localhost:8080"} 1024
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="http://localhost:8080",version="1.0"} 1
`)))
}
