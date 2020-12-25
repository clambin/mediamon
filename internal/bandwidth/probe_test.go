package bandwidth_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"

	"mediamon/internal/bandwidth"
	"mediamon/internal/metrics"
)

func tempFile(content []byte) (string, error) {
	filename := ""
	file, err := ioutil.TempFile("", "openvpn_")
	if err == nil {
		filename = file.Name()
		_, _ = file.Write(content)
		file.Close()
	}
	return filename, err
}

func TestProbe_Run(t *testing.T) {
	testCases := []struct {
		name        string
		content     []byte
		read, write float64
	}{
		{"empty", []byte(``), 0.0, 0.0},
		{"valid", []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,5624951995
TCP/UDP write bytes,2048
END`), 5624951995.0, 2048.0},
		{"invalid", []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,A
TCP/UDP write bytes,B
END`), 0.0, 0.0},
	}

	// valid/invalid file content

	for _, testCase := range testCases {

		filename, err := tempFile(testCase.content)

		if assert.Nil(t, err, testCase.name) {
			probe := bandwidth.NewProbe(filename)
			assert.NotNil(t, probe, testCase.name)

			probe.Run()

			read, ok := metrics.LoadValue("openvpn_client_tcp_udp_read_bytes_total")
			assert.True(t, ok, testCase.name)
			assert.Equal(t, testCase.read, read, testCase.name)

			write, ok := metrics.LoadValue("openvpn_client_tcp_udp_write_bytes_total")
			assert.True(t, ok, testCase.name)
			assert.Equal(t, testCase.write, write, testCase.name)

			_ = os.Remove(filename)
		}

	}

	// missing file

	bandwidth.NewProbe("invalidfile.txt").Run()

	read, ok := metrics.LoadValue("openvpn_client_tcp_udp_read_bytes_total")
	assert.True(t, ok)
	assert.Equal(t, 0.0, read)
	write, ok := metrics.LoadValue("openvpn_client_tcp_udp_write_bytes_total")
	assert.True(t, ok)
	assert.Equal(t, 0.0, write)
}
