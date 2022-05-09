package connectivity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clambin/mediamon/cache"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
	cache.Cache[bool]
	URL    string
	token  string
	client *http.Client
}

const httpTimeout = 10 * time.Second

// NewCollector creates a new Collector
func NewCollector(token string, proxyURL *url.URL, interval time.Duration) prometheus.Collector {
	transport := &http.Transport{}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	c := &Collector{
		token: token,
		client: &http.Client{
			Transport: transport,
			Timeout:   httpTimeout,
		},
	}
	c.Cache = cache.Cache[bool]{
		Duration: interval,
		Updater:  c.getState,
	}

	return c
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	value := 0.0
	if coll.Update() == true {
		value = 1.0
	}
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, value)
}

func (coll *Collector) getState() (state bool, err error) {
	err = coll.ping()
	up := err == nil

	return up, nil
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
	resp, err = coll.client.Do(req)

	if err == nil {
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

		if resp.StatusCode == http.StatusOK {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)

			if err == nil {
				decoder := json.NewDecoder(bytes.NewReader(body))
				err = decoder.Decode(&response)
			}
		} else {
			err = fmt.Errorf("%s", resp.Status)
		}
	}

	log.WithError(err).Debug("ping")

	return
}
