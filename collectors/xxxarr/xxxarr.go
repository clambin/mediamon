package xxxarr

import (
	"context"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/collectors/xxxarr/scraper"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

// Collector presents Sonarr/Radarr statistics as Prometheus metrics
type Collector struct {
	scraper.Scraper
	application string
	metrics     map[string]*prometheus.Desc
	transport   *httpclient.RoundTripper
	logger      *slog.Logger
}

// Config to create a collector
type Config struct {
	URL    string
	APIKey string
}

var (
	radarrCacheTable = httpclient.CacheTable{
		{Path: `/api/v3/system/status`, Expiry: time.Minute},
		{Path: `/api/v3/calendar`, Expiry: time.Minute},
		{Path: `/api/v3/movie`},
		{Path: `/api/v3/movie/[\d+]`, IsRegExp: true},
	}

	sonarrCacheTable = httpclient.CacheTable{
		{Path: `/api/v3/system/status`, Expiry: time.Minute},
		{Path: `/api/v3/calendar`, Expiry: time.Minute},
		{Path: `/api/v3/series`},
		{Path: `/api/v3/series/[\d+]`, IsRegExp: true},
		{Path: `/api/v3/episode/[\d+]`, IsRegExp: true},
	}
)

const (
	cacheExpiry     = 15 * time.Minute
	cleanupInterval = 5 * time.Minute
)

// NewRadarrCollector creates a new RadarrCollector
func NewRadarrCollector(url, apiKey string) *Collector {
	r := httpclient.NewRoundTripper(
		httpclient.WithCache{
			Table:           radarrCacheTable,
			DefaultExpiry:   cacheExpiry,
			CleanupInterval: cleanupInterval,
		},
		httpclient.WithRoundTripperMetrics{Namespace: "mediamon", Application: "radarr"},
	)

	return &Collector{
		Scraper: &scraper.RadarrScraper{
			Client: xxxarr.RadarrClient{
				Client: &http.Client{Transport: r},
				URL:    url,
				APIKey: apiKey,
			},
		},
		application: "radarr",
		metrics:     createMetrics("radarr", url),
		transport:   r,
		logger:      slog.Default().With(slog.String("application", "radarr")),
	}
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string) *Collector {
	r := httpclient.NewRoundTripper(
		httpclient.WithCache{
			Table:           sonarrCacheTable,
			DefaultExpiry:   cacheExpiry,
			CleanupInterval: cleanupInterval,
		},
		httpclient.WithRoundTripperMetrics{Namespace: "mediamon", Application: "sonarr"},
	)

	return &Collector{
		Scraper: &scraper.SonarrScraper{
			Client: xxxarr.SonarrClient{
				Client: &http.Client{Transport: r},
				URL:    url,
				APIKey: apiKey,
			},
		},
		application: "sonarr",
		metrics:     createMetrics("sonarr", url),
		transport:   r,
		logger:      slog.Default().With(slog.String("application", "sonarr")),
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range coll.metrics {
		ch <- metric
	}
	coll.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	stats, err := coll.Scraper.Scrape(context.Background())
	if err != nil {
		/*
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("mediamon_error",
					"Error getting "+coll.application+" metrics", nil, nil),
				err)
		*/
		coll.logger.Error("failed to collect `%s` metrics", err)
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
		// TODO: this fails if the same queued.Name is being downloaded multiple times. /metrics reports:
		// An error has occurred while serving metrics:
		//
		//2 error(s) occurred:
		//* collected metric "mediamon_xxxarr_queued_total_bytes" { label:<name:"application" value:"sonarr" > label:<name:"title" value:"name" > label:<name:"url" value:"http://sonarr:8989" > gauge:<value:6.476722065e+09 > } was collected before with the same name and label values
		//* collected metric "mediamon_xxxarr_queued_downloaded_bytes" { label:<name:"application" value:"sonarr" > label:<name:"title" value:"name" > label:<name:"url" value:"http://sonarr:8989" > gauge:<value:3.4451345e+07 > } was collected before with the same name and label values
		ch <- prometheus.MustNewConstMetric(coll.metrics["queued_total"], prometheus.GaugeValue, queued.TotalBytes, queued.Name)
		ch <- prometheus.MustNewConstMetric(coll.metrics["queued_downloaded"], prometheus.GaugeValue, queued.DownloadedBytes, queued.Name)
	}
	ch <- prometheus.MustNewConstMetric(coll.metrics["monitored"], prometheus.GaugeValue, float64(stats.Monitored))
	ch <- prometheus.MustNewConstMetric(coll.metrics["unmonitored"], prometheus.GaugeValue, float64(stats.Unmonitored))
	coll.transport.Collect(ch)

	slog.Debug(coll.application+" stats collected", "duration", time.Since(start))
}
