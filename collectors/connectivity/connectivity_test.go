package connectivity_test

import (
	"github.com/clambin/httpclient"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	m := httpclient.NewMetrics("foo", "")
	proxy, _ := url.Parse("http://localhost:8888")
	c := connectivity.NewCollector("123", proxy, 5*time.Minute, m)
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

	m := httpclient.NewMetrics("foo", "")
	c := connectivity.NewCollector("foo", nil, 5*time.Minute, m)
	c.URL = testServer.URL

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 1
`)))
}

func TestCollector_Collect_Down(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(down))
	defer testServer.Close()

	m := httpclient.NewMetrics("foo", "")
	c := connectivity.NewCollector("foo", nil, 5*time.Minute, m)
	c.URL = testServer.URL
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 0
`)))
}

func up(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`{}`))
}

func down(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "we're not home", http.StatusInternalServerError)
}
