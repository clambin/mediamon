package connectivity

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	tp := fakeTransport{pass: true}
	c := NewCollector(&http.Client{Transport: &tp}, 0)

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 1
`)))

	tp.pass = false
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 0
`)))
}

var _ http.RoundTripper = (*fakeTransport)(nil)

type fakeTransport struct {
	pass bool
}

func (f fakeTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	if !f.pass {
		return nil, errors.New("failed")
	}
	return &http.Response{StatusCode: http.StatusNoContent}, nil
}
