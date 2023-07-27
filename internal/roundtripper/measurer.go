package roundtripper

import (
	"github.com/clambin/go-common/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"time"
)

var _ httpclient.RequestMeasurer = RequestMeasurer{}

type RequestMeasurer struct {
	latency *prometheus.SummaryVec // measures latency of an API call
	errors  *prometheus.CounterVec // measures any errors returned by an API call
}

func NewRequestMeasurer(namespace, subsystem, application string) httpclient.RequestMeasurer {
	return RequestMeasurer{
		latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "api_latency"),
			Help:        "latency of Reporter API calls",
			ConstLabels: map[string]string{"application": application},
		}, []string{"method", "path"}),
		errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "api_errors_total"),
			Help:        "Number of failed Reporter API calls",
			ConstLabels: map[string]string{"application": application},
		}, []string{"method", "path"}),
	}
}

var idEliminator = regexp.MustCompile("^(?P<path>.+)/[0-9]+$")

func (r RequestMeasurer) MeasureRequest(req *http.Request, _ *http.Response, err error, duration time.Duration) {
	path := req.URL.Path
	if matches := idEliminator.FindAllStringSubmatch(path, 1); len(matches) == 1 && len(matches[0]) == 2 {
		path = matches[0][1]
	}
	r.latency.WithLabelValues(req.Method, path).Observe(duration.Seconds())
	var val float64
	if err != nil {
		val = 1
	}
	r.errors.WithLabelValues(req.Method, path).Add(val)
}

func (r RequestMeasurer) Describe(descs chan<- *prometheus.Desc) {
	r.latency.Describe(descs)
	r.errors.Describe(descs)
}

func (r RequestMeasurer) Collect(metrics chan<- prometheus.Metric) {
	r.latency.Collect(metrics)
	r.errors.Collect(metrics)
}
