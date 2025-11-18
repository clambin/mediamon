package plex

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestVersionCollector_Collect(t *testing.T) {
	c := newVersionCollector(
		fakeGetter{identity: plex.Identity{Version: "1.2.3"}},
		"http://localhost:8080",
		slog.New(slog.DiscardHandler),
	)

	expected := bytes.NewBufferString(`
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="http://localhost:8080",version="1.2.3"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(c, expected))
}
