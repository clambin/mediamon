package bandwidth_test

import (
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := bandwidth.NewCollector("", 5*time.Minute)
	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{"openvpn_client_tcp_udp_read_bytes_total", "openvpn_client_tcp_udp_write_bytes_total"} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name        string
		content     []byte
		pass        bool
		read, write float64
	}{
		{
			name:    "empty",
			content: []byte(``),
			pass:    false,
		},
		{
			name: "valid",
			content: []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,5624951995
TCP/UDP write bytes,2048
END`),
			pass:  true,
			read:  5624951995,
			write: 2048,
		},
		{
			name: "invalid",
			content: []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,A
TCP/UDP write bytes,B
END`),
			pass: false,
		},
	}

	// valid/invalid file content

	for _, testCase := range testCases {
		if filename, err := tempFile(testCase.content); assert.Nil(t, err, testCase.name) {
			c := bandwidth.NewCollector(filename, 5*time.Minute)
			ch := make(chan prometheus.Metric)
			go c.Collect(ch)

			if testCase.pass {
				read := <-ch
				assert.Equal(t, testCase.read, metrics.MetricValue(read).GetGauge().GetValue())
				write := <-ch
				assert.Equal(t, testCase.write, metrics.MetricValue(write).GetGauge().GetValue())
			} else {
				metric := <-ch
				assert.Equal(t, `Desc{fqName: "mediamon_error", help: "Error getting bandwidth statistics", constLabels: {}, variableLabels: []}`, metric.Desc().String())
			}
		}
	}
}

func tempFile(content []byte) (string, error) {
	filename := ""
	file, err := ioutil.TempFile("", "openvpn_")
	if err == nil {
		filename = file.Name()
		_, _ = file.Write(content)
		_ = file.Close()
	}
	return filename, err
}

func TestCollector_Collect_Failure(t *testing.T) {
	c := bandwidth.NewCollector("invalid file", 5*time.Minute)
	ch := make(chan prometheus.Metric)

	go c.Collect(ch)
	metric := <-ch
	assert.Equal(t, `Desc{fqName: "mediamon_error", help: "Error getting bandwidth statistics", constLabels: {}, variableLabels: []}`, metric.Desc().String())
}
