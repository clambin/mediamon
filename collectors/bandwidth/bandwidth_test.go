package bandwidth_test

import (
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestCollector_Describe(t *testing.T) {
	c := bandwidth.NewCollector("")
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{"openvpn_client_tcp_udp_read_bytes_total", "openvpn_client_tcp_udp_write_bytes_total"} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name    string
		content []byte
		pass    bool
		output  string
	}{
		{
			name: "valid",
			content: []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,5624951995
TCP/UDP write bytes,2048
END`),
			pass: true,
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
			name: "invalid",
			content: []byte(`OpenVPN STATISTICS
			Updated,Fri Dec 18 11:24:01 2020
			TCP/UDP read bytes,A
			TCP/UDP write bytes,B
			END
`),
			pass: true,
		},
	}

	for _, testCase := range testCases {
		filename, err := tempFile(testCase.content)
		require.NoError(t, err)

		c := bandwidth.NewCollector(filename)
		if testCase.pass {
			assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(testCase.output)))
		} else {
			err = testutil.CollectAndCompare(c, strings.NewReader(""))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), `Desc{fqName: "mediamon_error", help: "Error getting bandwidth statistics", constLabels: {}, variableLabels: []}`)
		}
		_ = os.Remove(filename)
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

/*
func TestCollector_Collect_Failure(t *testing.T) {
	c := bandwidth.NewCollector("invalid file")
	ch := make(chan prometheus.Metric)

	go c.Collect(ch)
	metric := <-ch
	assert.Equal(t, `Desc{fqName: "mediamon_error", help: "Error getting bandwidth statistics", constLabels: {}, variableLabels: []}`, metric.Desc().String())
}
*/
