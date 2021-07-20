package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestDuration collects request time statistics
	RequestDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "mediamon_request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})
)
