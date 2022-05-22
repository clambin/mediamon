package metrics

import (
	"github.com/clambin/go-metrics/client"
)

var (
	callerMetrics = client.NewMetrics("mediamon", "")

	// Latency collects request time statistics
	Latency = callerMetrics.Latency
	// Errors collects request error statistics
	Errors = callerMetrics.Errors
)
