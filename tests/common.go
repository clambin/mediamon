package tests

import (
	"github.com/clambin/gotools/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// ValidateMetric checks that a prometheus metric has a specified value
func ValidateMetric(metric prometheus.Metric, value float64, labelName, labelValue string) bool {
	if metrics.MetricValue(metric).GetGauge().GetValue() != value {
		return false
	}

	if labelName != "" {
		return metrics.MetricLabel(metric, labelName) == labelValue
	}

	return true
}
