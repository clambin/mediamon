package prowlarr

import (
	"context"
	"github.com/clambin/mediaclients/prowlarr"
	"github.com/clambin/mediamon/v2/internal/collectors/prowlarr/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strings"
	"testing"
)

func TestCollector(t *testing.T) {
	p := mocks.NewProwlarrClient(t)
	p.EXPECT().
		GetApiV1IndexerstatsWithResponse(context.Background(), (*prowlarr.GetApiV1IndexerstatsParams)(nil)).
		Return(&prowlarr.GetApiV1IndexerstatsResponse{JSON200: &prowlarr.IndexerStatsResource{
			Indexers: &[]prowlarr.IndexerStatistics{{
				IndexerId:             constP[int32](1),
				IndexerName:           constP("foo"),
				AverageResponseTime:   constP[int32](100),
				NumberOfQueries:       constP[int32](10),
				NumberOfFailedQueries: constP[int32](1),
				NumberOfGrabs:         constP[int32](2),
				NumberOfFailedGrabs:   constP[int32](1),
			}},
			UserAgents: &[]prowlarr.UserAgentStatistics{{
				UserAgent:       constP("foo"),
				NumberOfQueries: constP[int32](10),
				NumberOfGrabs:   constP[int32](1),
			}}}}, nil).
		Once()

	c := New("http://localhost", "", slog.Default())
	c.Collector.(*Collector).ProwlarrClient = p

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP mediamon_prowlarr_indexer_failed_grab_total Total number of failed grabs from this indexer
# TYPE mediamon_prowlarr_indexer_failed_grab_total counter
mediamon_prowlarr_indexer_failed_grab_total{application="prowlarr",indexer="foo",url="http://localhost"} 1

# HELP mediamon_prowlarr_indexer_failed_query_total Total number of failed queries to this indexer
# TYPE mediamon_prowlarr_indexer_failed_query_total counter
mediamon_prowlarr_indexer_failed_query_total{application="prowlarr",indexer="foo",url="http://localhost"} 1

# HELP mediamon_prowlarr_indexer_grab_total Total number of grabs from this indexer
# TYPE mediamon_prowlarr_indexer_grab_total counter
mediamon_prowlarr_indexer_grab_total{application="prowlarr",indexer="foo",url="http://localhost"} 2

# HELP mediamon_prowlarr_indexer_query_total Total number of queries to this indexer
# TYPE mediamon_prowlarr_indexer_query_total counter
mediamon_prowlarr_indexer_query_total{application="prowlarr",indexer="foo",url="http://localhost"} 10

# HELP mediamon_prowlarr_indexer_response_time Average response time in seconds
# TYPE mediamon_prowlarr_indexer_response_time gauge
mediamon_prowlarr_indexer_response_time{application="prowlarr",indexer="foo",url="http://localhost"} 0.1

# HELP mediamon_prowlarr_user_agent_grab_total Total number of grabs by user agent
# TYPE mediamon_prowlarr_user_agent_grab_total counter
mediamon_prowlarr_user_agent_grab_total{application="prowlarr",url="http://localhost",user_agent="foo"} 1

# HELP mediamon_prowlarr_user_agent_query_total Total number of queries by user agent
# TYPE mediamon_prowlarr_user_agent_query_total counter
mediamon_prowlarr_user_agent_query_total{application="prowlarr",url="http://localhost",user_agent="foo"} 10
`),
		"mediamon_prowlarr_indexer_grab_total",
		"mediamon_prowlarr_indexer_query_total",
		"mediamon_prowlarr_indexer_failed_grab_total",
		"mediamon_prowlarr_indexer_failed_query_total",
		"mediamon_prowlarr_indexer_response_time",
		"mediamon_prowlarr_user_agent_query_total",
		"mediamon_prowlarr_user_agent_grab_total",
	))
}

func constP[T any](t T) *T {
	return &t
}
