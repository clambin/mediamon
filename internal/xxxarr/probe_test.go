package xxxarr_test

import (
	"errors"
	"testing"

	"github.com/clambin/gotools/metrics"
	"github.com/stretchr/testify/assert"

	"mediamon/internal/xxxarr"
)

type client struct {
	application string
	fail        bool
}

func (client *client) GetApplication() string {
	return client.application
}

func (client *client) GetVersion() (string, error) {
	if client.fail {
		return "", errors.New("failing")
	}
	return "foo", nil
}

func (client *client) GetCalendar() (int, error) {
	if client.fail {
		return 0, errors.New("failing")
	}
	return 1, nil
}

func (client *client) GetQueue() (int, error) {
	if client.fail {
		return 0, errors.New("failing")
	}
	return 2, nil
}

func (client *client) GetMonitored() (int, int, error) {
	if client.fail {
		return 0, 0, errors.New("failing")
	}
	return 2, 1, nil
}

func TestProbe_Run(t *testing.T) {
	for _, application := range []string{"sonarr", "radarr"} {
		probe := xxxarr.Probe{XXXArrAPI: &client{application: application}}

		err := probe.Run()
		assert.Nil(t, err)

		value, _ := metrics.LoadValue("mediaserver_server_info", application, "foo")
		assert.Equal(t, float64(1), value)
		count, _ := metrics.LoadValue("mediaserver_calendar_count", application)
		assert.Equal(t, float64(1), count)
		count, _ = metrics.LoadValue("mediaserver_queued_count", application)
		assert.Equal(t, float64(2), count)
		count, _ = metrics.LoadValue("mediaserver_monitored_count", application)
		assert.Equal(t, float64(2), count)
		count, _ = metrics.LoadValue("mediaserver_unmonitored_count", application)
		assert.Equal(t, float64(1), count)
	}
}

func TestProbe_Fail(t *testing.T) {
	probe := xxxarr.Probe{&client{fail: true}}

	err := probe.Run()
	assert.NotNil(t, err)
}
