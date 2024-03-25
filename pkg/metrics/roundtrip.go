package metrics

import (
	"github.com/clambin/go-common/http/metrics"
	"net/http"
	"time"
)

var _ metrics.RequestMetrics = CustomizedRoundTripMetrics{}

type CustomizedRoundTripMetrics struct {
	metrics.RequestMetrics
	customize RequestCustomizer
}

type RequestCustomizer func(r *http.Request) *http.Request

func NewCustomizedRoundTripMetrics(namespace, subsystem string, labels map[string]string, f RequestCustomizer) metrics.RequestMetrics {
	return CustomizedRoundTripMetrics{
		RequestMetrics: metrics.NewRequestSummaryMetrics(namespace, subsystem, labels),
		customize:      f,
	}
}

func (m CustomizedRoundTripMetrics) Measure(req *http.Request, statusCode int, duration time.Duration) {
	m.RequestMetrics.Measure(m.customize(req), statusCode, duration)
}
