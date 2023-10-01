package connectivity

import (
	"fmt"
	"github.com/clambin/go-common/httpclient"
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
	HTTPClient *http.Client
	URL        string
	token      string
	transport  *httpclient.RoundTripper
	logger     *slog.Logger
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
func NewCollector(token string, proxyURL *url.URL, expiry time.Duration) *Collector {
	options := []httpclient.Option{
		httpclient.WithInstrumentedCache(httpclient.DefaultCacheTable, expiry, 2*expiry, "mediamon", "", "connectivity"),
		httpclient.WithMetrics("mediamon", "", "connectivity"),
	}
	if proxyURL != nil {
		options = append(options, httpclient.WithRoundTripper(&http.Transport{Proxy: http.ProxyURL(proxyURL)}))
	}

	r := httpclient.NewRoundTripper(options...)

	return &Collector{
		HTTPClient: &http.Client{
			Transport: r,
			Timeout:   httpTimeout,
		},
		token:     token,
		transport: r,
		logger:    slog.Default().With("collector", "connectivity"),
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
	c.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	var value float64
	if err := c.ping(); err == nil {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
	c.transport.Collect(ch)
	c.logger.Debug("stats collected", "duration", time.Since(start))
}

func (c *Collector) ping() error {
	URL := "https://ipinfo.io/"
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

	/*
		var response struct {
			IP       string
			Hostname string
			City     string
			Region   string
			Country  string
			Loc      string
			Org      string
			Postal   string
			Timezone string
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		return json.Unmarshal(body, &response)

	*/
	return nil
}
