package transmission

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	g := fakeTransmissionClient{
		sessionStats: transmission.SessionStats{
			ActiveTorrents: 1,
			PausedTorrents: 2,
			UploadRate:     25,
			DownloadRate:   100,
		},
		session: transmission.Session{
			Version: "foo",
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

var _ TransmissionClient = &fakeTransmissionClient{}

type fakeTransmissionClient struct {
	sessionStats transmission.SessionStats
	session      transmission.Session
	err          error
}

func (f fakeTransmissionClient) GetSession(_ context.Context, _ ...transmission.SessionField) (*transmission.Session, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &f.session, nil
}

func (f fakeTransmissionClient) GetSessionStats(_ context.Context) (*transmission.SessionStats, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &f.sessionStats, nil
}
