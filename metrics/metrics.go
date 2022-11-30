package metrics

import (
	"github.com/clambin/httpclient"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ClientMetrics = httpclient.NewMetrics("mediamon", "", prometheus.DefaultRegisterer)
)
