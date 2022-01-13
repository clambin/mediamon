package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Latency collects request time statistics
	Latency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "mediamon_request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	// Errors collects request error statistics
	Errors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mediamon_request_errors_total",
		Help: "API requests errors",
	}, []string{"application", "request"})
)
