package metrics_test

import (
	"github.com/clambin/mediamon/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCustomizedCacheMetrics(t *testing.T) {
	m := metrics.NewCustomizedCacheMetrics("", "", "", func(r *http.Request) *http.Request {
		return &http.Request{Method: r.Method, URL: &url.URL{Path: "/<redacted>"}}
	})

	req := http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/foo"}}

	m.Measure(&req, false)
	m.Measure(&req, true)

	assert.NoError(t, testutil.CollectAndCompare(m, strings.NewReader(`
# HELP http_cache_hit_total Number of times the cache was used
# TYPE http_cache_hit_total counter
http_cache_hit_total{method="POST",path="/<redacted>"} 1

# HELP http_cache_total Number of times the cache was consulted
# TYPE http_cache_total counter
http_cache_total{method="POST",path="/<redacted>"} 2
`)))
}
