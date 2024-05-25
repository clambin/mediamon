package connectivity_test

import (
	"errors"
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity"
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCollector_Collect(t *testing.T) {
	u, _ := url.Parse("http://localhost:8888")
	c := connectivity.NewCollector(u, 5*time.Minute, slog.Default())
	l := mocks.NewLocator(t)
	c.Locator = l

	l.EXPECT().Locate("").Return(0, 0, nil).Once()

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 1
`)))

	l.EXPECT().Locate("").Return(0, 0, errors.New("fail")).Once()

	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(`
# HELP openvpn_client_status OpenVPN client status
# TYPE openvpn_client_status gauge
openvpn_client_status 0
`)))

}
