package plex

import (
	"github.com/clambin/go-common/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

var _ httpclient.RequestMeasurer = Measurer{}

type Measurer struct {
	latency *prometheus.SummaryVec // measures latency of an API call
	errors  *prometheus.CounterVec // measures any errors returned by an API call
}

func newMeasurer(namespace, subsystem, application string) Measurer {
	return Measurer{
		latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "api_latency"),
			Help:        "latency of HTTP calls",
			ConstLabels: map[string]string{"application": application},
		}, []string{"method", "path"}),
		errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "api_errors_total"),
			Help:        "Number of failed HTTP calls",
			ConstLabels: map[string]string{"application": application},
		}, []string{"method", "path"}),
	}
}

func (m Measurer) MeasureRequest(req *http.Request, _ *http.Response, err error, duration time.Duration) {
	path := commonPath(req.URL.Path)
	m.latency.WithLabelValues(req.Method, path).Observe(duration.Seconds())
	var val float64
	if err != nil {
		val = 1
	}
	m.errors.WithLabelValues(req.Method, path).Add(val)
}

func commonPath(path string) string {
	for _, prefix := range []string{"/library/metadata", "/library/sections"} {
		if strings.HasPrefix(path, prefix) {
			return prefix
		}
	}
	return path
}

func (m Measurer) Describe(ch chan<- *prometheus.Desc) {
	m.latency.Describe(ch)
	m.errors.Describe(ch)
}

func (m Measurer) Collect(ch chan<- prometheus.Metric) {
	m.latency.Collect(ch)
	m.errors.Collect(ch)
}
