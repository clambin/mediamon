package connectivity_test

import (
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/tests"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := connectivity.NewCollector("123", "http://localhost:8888", 5*time.Minute)
	metrics := make(chan *prometheus.Desc)
	go c.Describe(metrics)

	for _, metricName := range []string{"openvpn_client_status"} {
		metric := <-metrics
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect_Up(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(up))
	defer testServer.Close()

	c := connectivity.NewCollector("foo", "", 5*time.Minute)
	c.(*connectivity.Collector).URL = testServer.URL
	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	metric := <-metrics
	assert.True(t, tests.ValidateMetric(metric, 1.0, "", ""))
}

func TestCollector_Collect_Down(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(down))
	defer testServer.Close()

	c := connectivity.NewCollector("foo", "", 5*time.Minute)
	c.(*connectivity.Collector).URL = testServer.URL
	metrics := make(chan prometheus.Metric)
	go c.Collect(metrics)

	metric := <-metrics
	assert.True(t, tests.ValidateMetric(metric, 0.0, "", ""))
}

func up(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`{}`))
}

func down(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "we're not home", http.StatusInternalServerError)
}
