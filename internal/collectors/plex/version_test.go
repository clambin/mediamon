package plex

import (
	"bytes"
	"github.com/clambin/mediaclients/plex"
	collectorbreaker "github.com/clambin/mediamon/v2/collector-breaker"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
)

func TestVersionCollector_Collect(t *testing.T) {
	p := mocks.NewGetter(t)
	p.EXPECT().GetIdentity(mock.Anything).Return(plex.Identity{Version: "1.2.3"}, nil)
	c := versionCollector{
		identityGetter: p,
		url:            "http://localhost:8080",
		logger:         slog.Default(),
	}

	expected := bytes.NewBufferString(`
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="http://localhost:8080",version="1.2.3"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(
		collectorbreaker.PassThroughCollector{Collector: c},
		expected,
	))
}
