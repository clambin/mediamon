package transmission_test

import (
	"context"
	"errors"
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/mediamon/internal/transmission"
	"github.com/stretchr/testify/assert"
	"testing"
)

type server struct {
	fail bool
}

func (server *server) GetVersion() (string, error) {
	if server.fail {
		return "", errors.New("failed")
	}
	return "foo", nil
}

func (server *server) GetStats() (int, int, int, int, error) {
	if server.fail {
		return 0, 0, 0, 0, errors.New("failed")
	}
	return 1, 2, 100, 25, nil
}

func TestMockedAPI(t *testing.T) {
	client := server{}
	probe := transmission.Probe{TransmissionAPI: &client}

	err := probe.Run(context.Background())
	assert.Nil(t, err)
	value, _ := metrics.LoadValue("mediaserver_server_info", "transmission", "foo")
	assert.Equal(t, float64(1), value)
	value, _ = metrics.LoadValue("mediaserver_active_torrent_count")
	assert.Equal(t, float64(1), value)
	value, _ = metrics.LoadValue("mediaserver_paused_torrent_count")
	assert.Equal(t, float64(2), value)
	value, _ = metrics.LoadValue("mediaserver_download_speed")
	assert.Equal(t, float64(100), value)
	value, _ = metrics.LoadValue("mediaserver_upload_speed")
	assert.Equal(t, float64(25), value)

	client.fail = true

	err = probe.Run(context.Background())
	assert.NotNil(t, err)
}
