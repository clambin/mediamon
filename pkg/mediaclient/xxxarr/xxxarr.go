package xxxarr

import (
	"github.com/clambin/go-metrics"
)

// Options contains options to alter Client behaviour
type Options struct {
	PrometheusMetrics metrics.APIClientMetrics // Prometheus metric to record API performance metrics
}
