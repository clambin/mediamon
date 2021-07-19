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

type Collector struct {
	cache.Cache
	URL    string
	token  string
	client *http.Client
	up     *prometheus.Desc
}

func NewCollector(token, proxyURL string, interval time.Duration) prometheus.Collector {
	c := &Collector{
		token:  token,
		client: getClient(proxyURL),
		up: prometheus.NewDesc(
			prometheus.BuildFQName("openvpn", "client", "status"),
			"OpenVPN client status",
			nil,
			nil,
		),
	}
	c.Cache = *cache.New(interval, false, c.getState)

	return c
}

func getClient(proxyURL string) (client *http.Client) {
	const httpTimeout = 10 * time.Second
	var proxy *url.URL

	if proxyURL != "" {
		var err error
		proxy, err = url.Parse(proxyURL)
		if err != nil {
			log.WithError(err).WithField("url", proxyURL).Warning("invalid proxy URL. ignoring")
			proxy = nil
		}
	}

	if proxy != nil {
		return &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
			Timeout:   httpTimeout,
		}
	}

	return &http.Client{
		Timeout: httpTimeout,
	}
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- coll.up
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	up := coll.Update().(bool)

	if up {
		ch <- prometheus.MustNewConstMetric(coll.up, prometheus.GaugeValue, 1.0)
	} else {
		ch <- prometheus.MustNewConstMetric(coll.up, prometheus.GaugeValue, 0.0)
	}
}

func (coll *Collector) getState() (state interface{}, err error) {
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
