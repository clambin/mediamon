package tests

import (
	"github.com/prometheus/client_golang/prometheus"
	pcg "github.com/prometheus/client_model/go"
)

// ValidateMetric checks that a prometheus metric has a specified value
func ValidateMetric(metric prometheus.Metric, value float64, labelName, labelValue string) bool {
	var m pcg.Metric
	if metric.Write(&m) == nil && *m.Gauge.Value == value {
		if labelName == "" {
			return true
		}

		for _, label := range m.Label {
			if *label.Name == labelName {
				if *label.Value == labelValue {
					return true
				}
			}
		}
	}

	return false
}
