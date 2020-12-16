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
