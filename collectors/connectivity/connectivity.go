package connectivity

import (
	"bytes"
	"encoding/json"
	"fmt"
	metrics2 "github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/caller"
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
	URL    string
	token  string
	Caller caller.Caller
}

const httpTimeout = 10 * time.Second

// NewCollector creates a new Collector
func NewCollector(token string, proxyURL *url.URL, interval time.Duration) prometheus.Collector {
	var httpClient *http.Client
	if proxyURL != nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: httpTimeout,
		}
	}

	return &Collector{
		token: token,
		Caller: caller.NewCacher(
			httpClient, "ipInfo",
			caller.Options{PrometheusMetrics: metrics2.APIClientMetrics{
				Latency: metrics.Latency,
				Errors:  metrics.Errors,
			}},
			[]caller.CacheTableEntry{},
			interval, 0,
		),
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	err := coll.ping()

	value := 0.0
	if err == nil {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
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
	resp, err = coll.Caller.Do(req)

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
