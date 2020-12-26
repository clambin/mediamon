package metrics_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"mediamon/internal/metrics"
)

func TestRun(t *testing.T) {
	assert.NotPanics(t, func() { metrics.Run(8080, true) })
	assert.Panics(t, func() { metrics.Run(8080, false) })
}

func TestPublish(t *testing.T) {
	// Unlabeled Gauge
	ok := metrics.Publish("upload_speed", 50)
	assert.True(t, ok)

	value, err := metrics.LoadValue("upload_speed")
	assert.Nil(t, err)
	assert.Equal(t, 50.0, value)

	// Labeled Gauge
	ok = metrics.Publish("plex_session_count", 5, "user")
	assert.True(t, ok)

	value, err = metrics.LoadValue("plex_session_count", "user")
	assert.Nil(t, err)
	assert.Equal(t, 5.0, value)

	// Invalid metric
	ok = metrics.Publish("not_a_metric", 5)
	assert.False(t, ok)

	_, err = metrics.LoadValue("not_a_metric")
	assert.NotNil(t, err)
}
