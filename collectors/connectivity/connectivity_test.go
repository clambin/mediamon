package connectivity_test

import (
	"github.com/clambin/go-metrics/tools"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	proxy, _ := url.Parse("http://localhost:8888")
	c := connectivity.NewCollector("123", proxy, 5*time.Minute)
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{"openvpn_client_status"} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect_Up(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(up))
	defer testServer.Close()

	c := connectivity.NewCollector("foo", nil, 5*time.Minute)
	c.(*connectivity.Collector).URL = testServer.URL
	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
	assert.Equal(t, 1.0, tools.MetricValue(metric).GetGauge().GetValue())
}

func TestCollector_Collect_Down(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(down))
	defer testServer.Close()

	c := connectivity.NewCollector("foo", nil, 5*time.Minute)
	c.(*connectivity.Collector).URL = testServer.URL
	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
	assert.Equal(t, 0.0, tools.MetricValue(metric).GetGauge().GetValue())
}

func up(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`{}`))
}

func down(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "we're not home", http.StatusInternalServerError)
}
