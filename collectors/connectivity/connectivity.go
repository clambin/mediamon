package connectivity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clambin/go-common/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"io"
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
	options := []httpclient.RoundTripperOption{
		httpclient.WithCache{
			DefaultExpiry:   expiry,
			CleanupInterval: 2 * expiry,
		},
		httpclient.WithRoundTripperMetrics{Namespace: "mediamon", Application: "connectivity"},
	}
	if proxyURL != nil {
		options = append(options, httpclient.WithRoundTripper{
			RoundTripper: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		})
	}

	r := httpclient.NewRoundTripper(options...)
	return &Collector{
		HTTPClient: &http.Client{
			Transport: r,
			Timeout:   httpTimeout,
		},
		transport: r,
		token:     token,
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
	coll.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	err := coll.ping()

	value := 0.0
	if err == nil {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
	coll.transport.Collect(ch)
}

func (coll *Collector) ping() (err error) {
	URL := "https://ipinfo.io/"
	if coll.URL != "" {
		URL = coll.URL
	}
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("token", coll.token)
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = coll.HTTPClient.Do(req)

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

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

	var body []byte
	body, err = io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&response)

	return
}
