package prowlarr

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clambin/mediaclients/prowlarr"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := prowlarr.IndexerStatsResource{
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
			}},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	t.Cleanup(ts.Close)

	want := `
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
`
	want = strings.ReplaceAll(want, "url=\"http://localhost\"", "url=\""+ts.URL+"\"")

	c, err := New(ts.URL, "1234", http.DefaultClient, slog.Default())
	require.NoError(t, err)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(want),
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
