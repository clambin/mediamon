package transmission

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	g := fakeTransmissionClient{
		sessionStats: transmissionrpc.SessionStats{
			ActiveTorrentCount: 1,
			PausedTorrentCount: 2,
			UploadSpeed:        25,
			DownloadSpeed:      100,
		},
		sessionArgs: transmissionrpc.SessionArguments{
			Version: constP("foo"),
		},
	}

	c, _ := NewCollector(http.DefaultClient, "", slog.New(slog.DiscardHandler))
	c.(*Collector).transmissionClient = &g

	e := strings.NewReader(`
# HELP mediamon_transmission_active_torrent_count Number of active torrents
# TYPE mediamon_transmission_active_torrent_count gauge
mediamon_transmission_active_torrent_count{url=""} 1

# HELP mediamon_transmission_download_speed Transmission download speed in bytes / sec
# TYPE mediamon_transmission_download_speed gauge
mediamon_transmission_download_speed{url=""} 100

# HELP mediamon_transmission_paused_torrent_count Number of paused torrents
# TYPE mediamon_transmission_paused_torrent_count gauge
mediamon_transmission_paused_torrent_count{url=""} 2

# HELP mediamon_transmission_upload_speed Transmission upload speed in bytes / sec
# TYPE mediamon_transmission_upload_speed gauge
mediamon_transmission_upload_speed{url=""} 25

# HELP mediamon_transmission_version version info
# TYPE mediamon_transmission_version gauge
mediamon_transmission_version{url="",version="foo"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(c, e))

	g.err = assert.AnError
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader("")))
}

func constP[T any](t T) *T {
	return &t
}

var _ TransmissionClient = &fakeTransmissionClient{}

type fakeTransmissionClient struct {
	sessionStats transmissionrpc.SessionStats
	sessionArgs  transmissionrpc.SessionArguments
	err          error
}

func (f fakeTransmissionClient) SessionArgumentsGetAll(_ context.Context) (sessionArgs transmissionrpc.SessionArguments, err error) {
	return f.sessionArgs, f.err
}

func (f fakeTransmissionClient) SessionStats(_ context.Context) (stats transmissionrpc.SessionStats, err error) {
	return f.sessionStats, f.err
}
