package metrics

import (
	"github.com/clambin/httpclient"
)

var (
	callerMetrics = httpclient.NewMetrics("mediamon", "")

	// Latency collects request time statistics
	Latency = callerMetrics.Latency
	// Errors collects request error statistics
	Errors = callerMetrics.Errors
)
