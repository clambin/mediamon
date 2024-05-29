package prowlarr

import (
	"context"
	"github.com/clambin/mediaclients/xxxarr"
	"github.com/clambin/mediamon/v2/internal/collectors/prowlarr/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	prowlarr := mocks.NewClient(t)
	prowlarr.EXPECT().GetIndexStats(context.Background()).Return(xxxarr.ProwlarrIndexersStats{Indexers: []xxxarr.ProwlarrIndexerStats{
		{
			IndexerId:           1,
			IndexerName:         "foo",
			AverageResponseTime: xxxarr.ProwlarrResponseTime(100 * time.Millisecond),
			NumberOfQueries:     10,
			NumberOfGrabs:       1,
		},
	}}, nil)

	c := New("http://localhost", "", slog.Default())
	c.Collector.(*Collector).client = prowlarr

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP mediamon_prowlarr_grab_total Total number of grabs from this indexer
# TYPE mediamon_prowlarr_grab_total counter
mediamon_prowlarr_grab_total{application="prowlarr",indexer="foo"} 1

# HELP mediamon_prowlarr_query_total Total number of queries to this indexer
# TYPE mediamon_prowlarr_query_total counter
mediamon_prowlarr_query_total{application="prowlarr",indexer="foo"} 10

# HELP mediamon_prowlarr_response_time Average response time in seconds
# TYPE mediamon_prowlarr_response_time gauge
mediamon_prowlarr_response_time{application="prowlarr",indexer="foo"} 0.1
`), "mediamon_prowlarr_grab_total", "mediamon_prowlarr_query_total", "mediamon_prowlarr_response_time"))
}
