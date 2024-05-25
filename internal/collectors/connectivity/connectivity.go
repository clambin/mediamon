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
	Locate(string) (float64, float64, error)
}

// Collector tests VPN connectivity by checking connection to https://ipinfo.io through a
// configured proxy
type Collector struct {
	Locator
	tpMetrics    metrics.RequestMetrics
	cacheMetrics roundtripper.CacheMetrics
	logger       *slog.Logger
}

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	Proxy    string
	Token    string
	Interval time.Duration
}

const httpTimeout = 10 * time.Second

// NewCollector creates a new Collector
func NewCollector(proxyURL *url.URL, expiry time.Duration, logger *slog.Logger) *Collector {
	cacheMetrics := roundtripper.NewCacheMetrics("mediamon", "", "connectivity")
	tpMetrics := metrics.NewRequestSummaryMetrics("mediamon", "", map[string]string{"application": "connectivity"})
	options := []roundtripper.Option{
		roundtripper.WithInstrumentedCache(roundtripper.DefaultCacheTable, expiry, 2*expiry, cacheMetrics),
		roundtripper.WithRequestMetrics(tpMetrics),
	}
	if proxyURL != nil {
		options = append(options, roundtripper.WithRoundTripper(&http.Transport{Proxy: http.ProxyURL(proxyURL)}))
	}
	httpClient := http.Client{
		Transport: roundtripper.New(options...),
		Timeout:   httpTimeout,
	}

	return &Collector{
		Locator:      iplocator.New(&httpClient),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       logger,
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
	c.tpMetrics.Describe(ch)
	c.cacheMetrics.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var value float64
	if _, _, err := c.Locate(""); err == nil {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
	c.tpMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
}
