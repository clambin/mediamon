package xxxarr

import (
	"context"
	"encoding/json"
	"github.com/clambin/go-metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestApiClient_WithMetrics(t *testing.T) {
	latencyMetric := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "xxxarr_request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	errorMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "xxxarr_request_errors",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	s := httptest.NewServer(http.HandlerFunc(handler))
	c := apiClient{
		HTTPClient:  http.DefaultClient,
		URL:         s.URL,
		APIKey:      "1234",
		application: "foo",
		options: Options{
			PrometheusMetrics: metrics.APIClientMetrics{
				Latency: latencyMetric,
				Errors:  errorMetric,
			},
		},
	}

	var response testStruct
	err := c.Get(context.Background(), "/foo", &response)
	require.NoError(t, err)
	assert.Equal(t, "bar", response.Name)
	assert.Equal(t, 42, response.Age)

	// validate that a metric was recorded
	ch := make(chan prometheus.Metric)
	go latencyMetric.Collect(ch)

	desc := <-ch
	assert.Equal(t, uint64(1), metrics.MetricValue(desc).GetSummary().GetSampleCount())

	// shut down the server
	s.Close()

	err = c.Get(context.Background(), "/foo", &response)
	require.Error(t, err)

	ch = make(chan prometheus.Metric)
	go errorMetric.Collect(ch)

	desc = <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(desc).GetCounter().GetValue())
}

func TestApiClient_Failures(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	c := apiClient{
		HTTPClient:  http.DefaultClient,
		URL:         s.URL,
		APIKey:      "4321",
		application: "foo",
	}

	ctx := context.Background()
	var response testStruct
	err := c.Get(ctx, "/foo", &response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "403 Forbidden")

	c.APIKey = "1234"
	err = c.Get(ctx, "/bar", &response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")

	err = c.Get(ctx, "foo", &response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to create request")

	s.Close()
	err = c.Get(ctx, "/foo", &response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connect: connection refused")

}

func handler(w http.ResponseWriter, req *http.Request) {
	// check auth
	if req.Header.Get("X-Api-Key") != "1234" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if req.URL.Path != "/foo" {
		http.Error(w, "invalid endpoint", http.StatusNotFound)
		return
	}

	response := testStruct{
		Name: "bar",
		Age:  42,
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
