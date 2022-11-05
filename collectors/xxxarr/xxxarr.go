package xxxarr

import (
	"context"
	"github.com/clambin/httpclient"
	"github.com/clambin/mediamon/collectors/xxxarr/scraper"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"time"
)

// Collector presents Sonarr/Radarr statistics as Prometheus metrics
type Collector struct {
	scraper.Scraper
	application string
	metrics     map[string]*prometheus.Desc
}

// Config to create a collector
type Config struct {
	URL    string
	APIKey string
}

var (
	radarrCacheTable = []httpclient.CacheTableEntry{
		{Endpoint: `/api/v3/system/status`, Expiry: time.Minute},
		{Endpoint: `/api/v3/calendar`, Expiry: time.Minute},
		{Endpoint: `/api/v3/movie`},
		{Endpoint: `/api/v3/movie/[\d+]`, IsRegExp: true},
	}

	sonarrCacheTable = []httpclient.CacheTableEntry{
		{Endpoint: `/api/v3/system/status`, Expiry: time.Minute},
		{Endpoint: `/api/v3/calendar`, Expiry: time.Minute},
		{Endpoint: `/api/v3/series`},
		{Endpoint: `/api/v3/series/[\d+]`, IsRegExp: true},
		{Endpoint: `/api/v3/episode/[\d+]`, IsRegExp: true},
	}
)

const (
	cacheExpiry     = 15 * time.Minute
	cleanupInterval = 5 * time.Minute
)

// NewRadarrCollector creates a new RadarrCollector
func NewRadarrCollector(url, apiKey string) *Collector {
	options := httpclient.Options{PrometheusMetrics: httpclient.Metrics{
		Latency: metrics.Latency,
		Errors:  metrics.Errors,
	}}
	c := httpclient.NewCacher(nil, "radarr", options, radarrCacheTable, cacheExpiry, cleanupInterval)

	return &Collector{
		Scraper: &scraper.RadarrScraper{
			Client: xxxarr.NewRadarrClientWithCaller(apiKey, url, c),
		},
		application: "radarr",
		metrics:     createMetrics("radarr", url),
	}
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string) *Collector {
	options := httpclient.Options{PrometheusMetrics: httpclient.Metrics{
		Latency: metrics.Latency,
		Errors:  metrics.Errors,
	}}
	c := httpclient.NewCacher(nil, "sonarr", options, sonarrCacheTable, cacheExpiry, cleanupInterval)
	return &Collector{
		Scraper: &scraper.SonarrScraper{
			Client: xxxarr.NewSonarrClientWithCaller(apiKey, url, c),
		},
		application: "sonarr",
		metrics:     createMetrics("sonarr", url),
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range coll.metrics {
		ch <- metric
	}
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := coll.Scraper.Scrape(context.Background())
	if err != nil {
		/*
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("mediamon_error",
					"Error getting "+coll.application+" metrics", nil, nil),
				err)
		*/
		log.WithError(err).Warningf("failed to collect `%s` metrics", coll.application)
		return
	}

	ch <- prometheus.MustNewConstMetric(coll.metrics["version"], prometheus.GaugeValue, float64(1), stats.Version)
	for key, value := range stats.Health {
		ch <- prometheus.MustNewConstMetric(coll.metrics["health"], prometheus.GaugeValue, float64(value), key)
	}
	for _, title := range stats.Calendar {
		ch <- prometheus.MustNewConstMetric(coll.metrics["calendar"], prometheus.GaugeValue, 1.0, title)
	}

	ch <- prometheus.MustNewConstMetric(coll.metrics["queued_count"], prometheus.GaugeValue, float64(len(stats.Queued)))
	for _, queued := range stats.Queued {
		ch <- prometheus.MustNewConstMetric(coll.metrics["queued_total"], prometheus.GaugeValue, queued.TotalBytes, queued.Name)
		ch <- prometheus.MustNewConstMetric(coll.metrics["queued_downloaded"], prometheus.GaugeValue, queued.DownloadedBytes, queued.Name)
	}
	ch <- prometheus.MustNewConstMetric(coll.metrics["monitored"], prometheus.GaugeValue, float64(stats.Monitored))
	ch <- prometheus.MustNewConstMetric(coll.metrics["unmonitored"], prometheus.GaugeValue, float64(stats.Unmonitored))
}
