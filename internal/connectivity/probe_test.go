package connectivity_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/clambin/httpstub"
	"github.com/stretchr/testify/assert"

	"mediamon/internal/connectivity"
	"mediamon/internal/metrics"
)

func TestProbe_Run(t *testing.T) {
	connectivity.NewProbeWithHTTPClient(httpstub.NewTestClient(loopback), "").Run()

	value, ok := metrics.LoadValue("openvpn_client_status")
	assert.True(t, ok)
	assert.Equal(t, 1.0, value)
}

func TestProbe_Run_Fail(t *testing.T) {
	connectivity.NewProbeWithHTTPClient(httpstub.NewTestClient(httpstub.Failing), "").Run()

	value, ok := metrics.LoadValue("openvpn_client_status")
	assert.True(t, ok)
	assert.Equal(t, 0.0, value)
}

// lookup function

func loopback(_ *http.Request) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(ipinfoResponse)),
	}
}

// Response

const ipinfoResponse = `{
  "ip": "1.2.3.4",
  "hostname": "example.com",
  "city": "City",
  "region": "Region",
  "country": "BE",
  "loc": "Loc",
  "org": "Org",
  "postal": "1234",
  "timezone": "Europe/Brussels"
}`
