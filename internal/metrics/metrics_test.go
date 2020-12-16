package metrics_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"mediamon/internal/metrics"
)

func TestLoadValue(t *testing.T) {
	metrics.SaveValue("metric", 12, "label1", "label2")
	loaded, ok := metrics.LoadValue("metric", "label1", "label2")
	assert.True(t, ok)
	assert.Equal(t, float64(12), loaded)

	loaded, ok = metrics.LoadValue("metric", "label1", "label3")
	assert.False(t, ok)
}

func TestInit(t *testing.T) {
	assert.NotPanics(t, func() { metrics.Init(8080) })
	assert.Panics(t, func() { metrics.Init(8080) })
}

func TestPublish(t *testing.T) {
	// Unlabeled Gauge
	ok := metrics.Publish("upload_speed", 50)
	assert.True(t, ok)

	value, ok := metrics.LoadValue("upload_speed")
	assert.True(t, ok)
	assert.Equal(t, 50.0, value)

	// Labeled Gauge
	ok = metrics.Publish("plex_session_count", 5, "user")
	assert.True(t, ok)

	value, ok = metrics.LoadValue("plex_session_count", "user")
	assert.True(t, ok)
	assert.Equal(t, 5.0, value)

	// Invalid metric
	ok = metrics.Publish("not_a_metric", 5)
	assert.False(t, ok)

	_, ok = metrics.LoadValue("not_a_metric")
	assert.False(t, ok)
}
