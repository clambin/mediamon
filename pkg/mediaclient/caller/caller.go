package caller

import (
	"github.com/clambin/go-metrics"
	"net/http"
)

// Caller interface of a generic API caller
type Caller interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// Client implements the Caller interface. If provided by Options, it will collect performance metrics of the API calls
// and record them for Prometheus to scrape.
type Client struct {
	HTTPClient  *http.Client
	Options     Options
	Application string
}

var _ Caller = &Client{}

// Options contains options to alter Client behaviour
type Options struct {
	PrometheusMetrics metrics.APIClientMetrics // Prometheus metric to record API performance metrics
}

// Do implements the Caller's Do() method. It sends the request and records performance metrics of the call.
// Currently, it records the request's duration (i.e. latency) and error rate.
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	endpoint := req.URL.Path
	timer := c.Options.PrometheusMetrics.MakeLatencyTimer(c.Application, endpoint)

	resp, err = c.HTTPClient.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}
	c.Options.PrometheusMetrics.ReportErrors(err, c.Application, endpoint)
	return
}
