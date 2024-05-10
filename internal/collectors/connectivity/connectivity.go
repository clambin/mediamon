package connectivity

import (
	"fmt"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
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

// Collector tests VPN connectivity by checking connection to https://ipinfo.io through a
// configured proxy
type Collector struct {
	HTTPClient   *http.Client
	URL          string
	token        string
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
func NewCollector(token string, proxyURL *url.URL, expiry time.Duration, logger *slog.Logger) *Collector {
	cacheMetrics := roundtripper.NewCacheMetrics("mediamon", "", "connectivity")
	tpMetrics := metrics.NewRequestSummaryMetrics("mediamon", "", map[string]string{"application": "connectivity"})
	options := []roundtripper.Option{
		roundtripper.WithInstrumentedCache(roundtripper.DefaultCacheTable, expiry, 2*expiry, cacheMetrics),
		roundtripper.WithRequestMetrics(tpMetrics),
	}
	if proxyURL != nil {
		options = append(options, roundtripper.WithRoundTripper(&http.Transport{Proxy: http.ProxyURL(proxyURL)}))
	}

	return &Collector{
		HTTPClient: &http.Client{
			Transport: roundtripper.New(options...),
			Timeout:   httpTimeout,
		},
		token:        token,
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
	if err := c.ping(); err == nil {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
	c.tpMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
}

func (c *Collector) ping() error {
	URL := "https://ipinfo.io"
	if c.URL != "" {
		URL = c.URL
	}
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("token", c.token)
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}
