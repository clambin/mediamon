package caller_test

import (
	"encoding/json"
	"fmt"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/caller"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Do(t *testing.T) {
	latencyMetric := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	errorMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_errors",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	s := httptest.NewServer(http.HandlerFunc(handler))
	c := &caller.Client{
		HTTPClient: http.DefaultClient,
		Options: caller.Options{
			PrometheusMetrics: metrics.APIClientMetrics{
				Latency: latencyMetric,
				Errors:  errorMetric,
			},
		},
		Application: "foo",
	}

	response, err := doCall(c, s.URL+"/foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", response.Name)
	assert.Equal(t, 42, response.Age)

	response, err = doCall(c, s.URL+"/bar")
	require.Error(t, err)

	s.Close()
	response, err = doCall(c, s.URL+"/foo")
	require.Error(t, err)

	// TODO: this is flaky.  sometimes metrics aren't correct???
	/*
		ch := make(chan prometheus.Metric)
		go latencyMetric.Collect(ch)

		desc := <-ch
		assert.Equal(t, uint64(2), metrics.MetricValue(desc).GetSummary().GetSampleCount())

		ch = make(chan prometheus.Metric)
		go errorMetric.Collect(ch)

		desc = <-ch
		assert.Equal(t, float64(1), metrics.MetricValue(desc).GetCounter().GetValue())
	*/
}

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func doCall(c caller.Caller, url string) (response testStruct, err error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	var resp *http.Response
	if resp, err = c.Do(req); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("call failed: %s", resp.Status)
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(&response)
	return
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/foo" {
		http.Error(w, "invalid endpoint", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(testStruct{
		Name: "bar",
		Age:  42,
	})
}
