package metrics_test

import (
	"github.com/clambin/mediamon/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCustomizedRoundTripMetrics(t *testing.T) {
	m := metrics.NewCustomizedRoundTripMetrics("", "", "", func(r *http.Request) *http.Request {
		return &http.Request{Method: r.Method, URL: &url.URL{Path: "/<redacted>"}}
	})

	req := http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/foo"}}
	resp := http.Response{StatusCode: http.StatusOK}

	m.Measure(&req, &resp, nil, time.Second)

	assert.NoError(t, testutil.CollectAndCompare(m, strings.NewReader(`
# HELP http_request_duration_seconds http request duration in seconds
# TYPE http_request_duration_seconds summary
http_request_duration_seconds_sum{code="200",method="POST",path="/<redacted>"} 1
http_request_duration_seconds_count{code="200",method="POST",path="/<redacted>"} 1

# HELP http_requests_total total number of http requests
# TYPE http_requests_total counter
http_requests_total{code="200",method="POST",path="/<redacted>"} 1
`)))
}
