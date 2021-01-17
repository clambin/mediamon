package transmission_test

import (
	"errors"
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/mediamon/internal/transmission"
	"github.com/stretchr/testify/assert"
	"testing"
)

type client struct {
	fail bool
}

func (client *client) GetVersion() (string, error) {
	if client.fail {
		return "", errors.New("failed")
	}
	return "foo", nil
}

func (client *client) GetStats() (int, int, int, int, error) {
	if client.fail {
		return 0, 0, 0, 0, errors.New("failed")
	}
	return 1, 2, 100, 25, nil
}

func TestMockedAPI(t *testing.T) {
	client := client{}
	probe := transmission.Probe{TransmissionAPI: &client}

	err := probe.Run()
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

	err = probe.Run()
	assert.NotNil(t, err)
}
