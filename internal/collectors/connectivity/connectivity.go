package connectivity

import (
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediamon/v2/pkg/iplocator"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

var (
	upMetric = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "status"),
		"OpenVPN client status",
		nil,
		nil,
	)
)

type Locator interface {
	Locate(string) (iplocator.Location, error)
}

// Collector tests network connectivity by querying the IP address location through ip-api.com
type Collector struct {
	Locator
	requestMetrics metrics.RequestMetrics
	cacheMetrics   roundtripper.CacheMetrics
	logger         *slog.Logger
}

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	Proxy    string
	Token    string
	Interval time.Duration
}

const httpTimeout = 10 * time.Second

// NewCollector creates a new Collector. proxyURL should be the URL of the transmission openvpn proxy. If expiration is set,
// IP address location requests are cached for that amount of time.
func NewCollector(proxyURL *url.URL, expiration time.Duration, logger *slog.Logger) *Collector {
	cacheMetrics := roundtripper.NewCacheMetrics("mediamon", "", "connectivity")
	requestMetrics := metrics.NewRequestMetrics(metrics.Options{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "connectivity"},
	})

	options := make([]roundtripper.Option, 0, 3)
	if expiration > 0 {
		options = append(options, roundtripper.WithCache(roundtripper.CacheOptions{
			DefaultExpiration: expiration,
			CleanupInterval:   2 * expiration,
			CacheMetrics:      cacheMetrics,
		}))
	}
	options = append(options, roundtripper.WithRequestMetrics(requestMetrics))
	if proxyURL != nil {
		options = append(options, roundtripper.WithRoundTripper(&http.Transport{Proxy: http.ProxyURL(proxyURL)}))
	}
	httpClient := http.Client{
		Transport: roundtripper.New(options...),
		Timeout:   httpTimeout,
	}

	return &Collector{
		Locator:        iplocator.New(&httpClient),
		requestMetrics: requestMetrics,
		cacheMetrics:   cacheMetrics,
		logger:         logger,
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
	c.requestMetrics.Describe(ch)
	c.cacheMetrics.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var value float64
	if _, err := c.Locate(""); err == nil {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
	c.requestMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
}
