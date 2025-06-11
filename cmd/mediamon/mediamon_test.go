package main

import (
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	rootCmd.Version = "ci/cd"
	rootCmd.SetContext(t.Context())
	rootCmd.SetArgs([]string{
		"--config", "testdata/config.yaml",
		"--openvpn.bandwidth.filename", "testdata/client.status",
	})

	go func() { _ = rootCmd.Execute() }()

	assert.Eventually(t, func() bool {
		_, err := http.Get("http://127.0.0.1:9090/metrics")
		return err == nil
	}, 5*time.Second, time.Millisecond*100)

	assert.NoError(t, testutil.GatherAndCompare(
		prometheus.DefaultGatherer,
		strings.NewReader(`
# HELP openvpn_client_tcp_udp_read_bytes_total OpenVPN client bytes read
# TYPE openvpn_client_tcp_udp_read_bytes_total gauge
openvpn_client_tcp_udp_read_bytes_total 5.893220736e+09

# HELP openvpn_client_tcp_udp_write_bytes_total OpenVPN client bytes written
# TYPE openvpn_client_tcp_udp_write_bytes_total gauge
openvpn_client_tcp_udp_write_bytes_total 1.882796878e+09
`),
		"openvpn_client_tcp_udp_read_bytes_total",
		"openvpn_client_tcp_udp_write_bytes_total",
	))
}

func Test_createCollectors(t *testing.T) {
	v := viper.New()
	v.Set("transmission.url", "http://transmission:80")
	v.Set("sonarr.url", "http://sonarr:80")
	v.Set("radarr.url", "http://radarr:80")
	v.Set("plex.url", "http://plex:80")
	v.Set("openvpn.connectivity.proxy", "http://proxy:8080")
	v.Set("openvpn.bandwidth.filename", "/data/client.status")

	collectors := createCollectors("ci/cd", v, slog.New(slog.DiscardHandler))
	assert.Len(t, collectors, 6)
}

func Test_parseProxy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  assert.ErrorAssertionFunc
		expected string
	}{
		{
			name:     "valid proxy url",
			input:    "http://proxy:8080",
			wantErr:  assert.NoError,
			expected: "http://proxy:8080",
		},
		{
			name:     "valid proxy url (no port)",
			input:    "http://proxy",
			wantErr:  assert.NoError,
			expected: "http://proxy",
		},
		{
			name:    "invalid proxy url",
			input:   "proxy",
			wantErr: assert.Error,
		},
		{
			name:    "invalid proxy url",
			input:   "\001",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := parseProxy(tt.input)
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}
