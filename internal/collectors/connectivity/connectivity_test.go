package connectivity_test

import (
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCollector_Collect_Up(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(up))
	defer testServer.Close()

	c := connectivity.NewCollector("foo", nil, 5*time.Minute)
	c.URL = testServer.URL

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP mediamon_http_cache_total Number of times the cache was consulted
# TYPE mediamon_http_cache_total counter
mediamon_http_cache_total{application="connectivity",method="GET",path="/"} 1

# HELP mediamon_http_requests_total total number of http requests
# TYPE mediamon_http_requests_total counter
mediamon_http_requests_total{application="connectivity",code="200",method="GET",path="/"} 1

# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 1
`), "openvpn_client_status", "mediamon_http_cache_total", "mediamon_http_requests_total"))
}

func TestCollector_Collect_Down(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(down))
	defer testServer.Close()

	c := connectivity.NewCollector("foo", nil, 5*time.Minute)
	c.URL = testServer.URL
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP mediamon_http_cache_total Number of times the cache was consulted
# TYPE mediamon_http_cache_total counter
mediamon_http_cache_total{application="connectivity",method="GET",path="/"} 1

# HELP mediamon_http_requests_total total number of http requests
# TYPE mediamon_http_requests_total counter
mediamon_http_requests_total{application="connectivity",code="500",method="GET",path="/"} 1

# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 0
`), "openvpn_client_status", "mediamon_http_cache_total", "mediamon_http_requests_total"))
}

func up(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("token") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(`{}`))
}

func down(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "we're not home", http.StatusInternalServerError)
}
