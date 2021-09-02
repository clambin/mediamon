package tests

import (
	"github.com/clambin/gotools/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// ValidateMetric checks that a prometheus metric has a specified value
func ValidateMetric(metric prometheus.Metric, value float64, args ...string) bool {
	if metrics.MetricValue(metric).GetGauge().GetValue() != value {
		return false
	}

	if len(args) == 0 {
		return true
	}

	if len(args)%2 != 0 {
		panic("args must be a list of label/value pairs")
	}

	for i := 0; i < len(args); i += 2 {
		if metrics.MetricLabel(metric, args[i]) != args[i+1] {
			return false
		}
	}

	return true
}
