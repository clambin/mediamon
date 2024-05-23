package plex

import (
	"bytes"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/plex/mocks"
	"github.com/clambin/mediamon/v2/pkg/breaker"
	collector_breaker "github.com/clambin/mediamon/v2/pkg/collector-breaker"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
	"time"
)

func TestVersionCollector_Collect(t *testing.T) {
	p := mocks.NewGetter(t)
	p.EXPECT().GetIdentity(mock.Anything).Return(plex.Identity{Version: "1.2.3"}, nil)
	c := versionCollector{
		versionGetter: p,
		url:           "http://localhost:8080",
		logger:        slog.Default(),
	}
	cfg := breaker.Configuration{
		FailureThreshold: 1,
		OpenDuration:     time.Minute,
		SuccessThreshold: 1,
	}
	cb := collector_breaker.New(c, cfg, slog.Default())

	expected := bytes.NewBufferString(`
# HELP mediamon_plex_version version info
# TYPE mediamon_plex_version gauge
mediamon_plex_version{url="http://localhost:8080",version="1.2.3"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(cb, expected))
}
