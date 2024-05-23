package bandwidth_test

import (
	"github.com/clambin/mediamon/v2/internal/collectors/bandwidth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := bandwidth.NewCollector("", slog.Default())
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{"openvpn_client_tcp_udp_read_bytes_total", "openvpn_client_tcp_udp_write_bytes_total"} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		output  string
	}{
		{
			name: "valid",
			content: []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,5624951995
TCP/UDP write bytes,2048
END`),
			output: `
# HELP openvpn_client_tcp_udp_read_bytes_total OpenVPN client bytes read
# TYPE openvpn_client_tcp_udp_read_bytes_total gauge
openvpn_client_tcp_udp_read_bytes_total 5.624951995e+09
# HELP openvpn_client_tcp_udp_write_bytes_total OpenVPN client bytes written
# TYPE openvpn_client_tcp_udp_write_bytes_total gauge
openvpn_client_tcp_udp_write_bytes_total 2048
`,
		},
		{
			name: "invalid values",
			content: []byte(`OpenVPN STATISTICS
			Updated,Fri Dec 18 11:24:01 2020
			TCP/UDP read bytes,A
			TCP/UDP write bytes,B
			END
`),
		},
		{
			name: "incomplete file",
			content: []byte(`OpenVPN STATISTICS
			Updated,Fri Dec 18 11:24:01 2020
			TCP/UDP read bytes,1
			TCP/UDP read bytes,1
			END
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			filename, err := tempFile(tt.content)
			require.NoError(t, err)

			c := bandwidth.NewCollector(filename, slog.Default())
			assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(tt.output)))

			assert.NoError(t, os.Remove(filename))
		})
	}
}

func tempFile(content []byte) (string, error) {
	filename := ""
	file, err := os.CreateTemp("", "openvpn_")
	if err == nil {
		filename = file.Name()
		_, _ = file.Write(content)
		_ = file.Close()
	}
	return filename, err
}
