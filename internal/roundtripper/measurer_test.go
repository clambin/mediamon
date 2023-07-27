package roundtripper_test

import (
	"bytes"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/v2/internal/roundtripper"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRequestMetrics(t *testing.T) {
	r := httpclient.NewRoundTripper(
		httpclient.WithCustomMetrics(roundtripper.NewRequestMeasurer("foo", "bar", "snafu")),
		httpclient.WithRoundTripper(httpclient.RoundTripperFunc(func(request *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK}, nil
		})),
	)

	req, _ := http.NewRequest(http.MethodGet, "/test/123", nil)
	resp, err := r.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, 1, testutil.CollectAndCount(r, `foo_bar_api_latency`))
	assert.NoError(t, testutil.CollectAndCompare(r, bytes.NewBufferString(`
# HELP foo_bar_api_errors_total Number of failed Reporter API calls
# TYPE foo_bar_api_errors_total counter
foo_bar_api_errors_total{application="snafu",method="GET",path="/test"} 0
`), `foo_bar_api_errors_total`))
}
