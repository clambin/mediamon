package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics contains prometheus metrics to capture during API calls
type PrometheusMetrics struct {
	Latency *prometheus.SummaryVec
	Errors  *prometheus.CounterVec
}

func (pm PrometheusMetrics) ReportErrors(err error, labelValues ...string) {
	if err != nil && pm.Errors != nil {
		pm.Errors.WithLabelValues(labelValues...).Add(1.0)
	}
}

func (pm PrometheusMetrics) MakeLatencyTimer(labelValues ...string) (timer *prometheus.Timer) {
	if pm.Latency != nil {
		timer = prometheus.NewTimer(pm.Latency.WithLabelValues(labelValues...))
	}
	return
}
