package metrics_test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"mediamon/pkg/metrics"
	"testing"
)

func TestGauge(t *testing.T) {
	gauge := metrics.NewGauge(prometheus.GaugeOpts{
		Name: "test_gauge",
		Help: "Gauge test",
	})

	assert.NotNil(t, gauge)
	gauge.Set(1)

	value, err := metrics.LoadValue("test_gauge")
	assert.Nil(t, err)
	assert.Equal(t, float64(1), value)
}

func TestGaugeVec(t *testing.T) {
	gauge := metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_gaugevec",
		Help: "Gauge test",
	}, []string{"host"})

	assert.NotNil(t, gauge)
	gauge.WithLabelValues("host1").Set(1)

	value, err := metrics.LoadValue("test_gaugevec", "host1")
	assert.Nil(t, err)
	assert.Equal(t, float64(1), value)

}

func TestInvalid(t *testing.T) {
	_, err := metrics.LoadValue("invalid")
	assert.NotNil(t, err)
}
