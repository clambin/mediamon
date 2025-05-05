package bandwidth

import (
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	content := []byte(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,5624951995
TCP/UDP write bytes,2048
END`)
	want := `
# HELP openvpn_client_tcp_udp_read_bytes_total OpenVPN client bytes read
# TYPE openvpn_client_tcp_udp_read_bytes_total gauge
openvpn_client_tcp_udp_read_bytes_total 5.624951995e+09
# HELP openvpn_client_tcp_udp_write_bytes_total OpenVPN client bytes written
# TYPE openvpn_client_tcp_udp_write_bytes_total gauge
openvpn_client_tcp_udp_write_bytes_total 2048
`

	filename, err := tempFile(content)
	require.NoError(t, err)

	c := NewCollector(filename, slog.Default())
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(want)))
	assert.NoError(t, os.Remove(filename))
	assert.Error(t, testutil.CollectAndCompare(c, strings.NewReader(want)))
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

func TestCollector_readStats(t *testing.T) {
	type want struct {
		err   assert.ErrorAssertionFunc
		stats bandwidthStats
	}
	tests := []struct {
		name    string
		content string
		want    want
	}{
		{
			name: "valid",
			content: `OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,1024
TCP/UDP write bytes,2048
END`,
			want: want{err: assert.NoError, stats: bandwidthStats{read: 1024, written: 2048}},
		},
		{
			name:    "empty",
			content: ``,
			want:    want{err: assert.Error},
		},
		{
			name:    "invalid line",
			content: `TCP/UDP read bytes,1024,100`,
			want:    want{err: assert.Error},
		},
		{
			name:    "invalid value",
			content: `TCP/UDP read bytes,102A`,
			want:    want{err: assert.Error},
		},
		{
			name:    "unknown line",
			content: `foo`,
			want:    want{err: assert.Error},
		},
		{
			name: "no read bytes",
			content: `OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP write bytes,2048
END`,
			want: want{err: assert.Error},
		},
		{
			name: "no write bytes",
			content: `OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,1024
END`,
			want: want{err: assert.Error},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stats, err := readStats(strings.NewReader(tt.content))
			tt.want.err(t, err)
			assert.Equal(t, tt.want.stats, stats)
		})
	}
}
