package xxxarr

import (
	"context"
	"errors"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/v2/collectors/xxxarr/scraper"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"time"
)

// Collector presents Sonarr/Radarr statistics as Prometheus metrics
type Collector struct {
	Scraper
	application string
	metrics     map[string]*prometheus.Desc
	transport   *httpclient.RoundTripper
	logger      *slog.Logger
}

// Scraper provides a generic means of getting stats from Sonarr or Radarr
//
//go:generate mockery --name Scraper
type Scraper interface {
	Scrape(ctx context.Context) (scraper.Stats, error)
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
		httpclient.WithInstrumentedCache(radarrCacheTable, cacheExpiry, cleanupInterval, "mediamon", "", "radarr"),
		httpclient.WithMetrics("mediamon", "", "radarr"),
	)

	return &Collector{
		Scraper:     scraper.RadarrScraper{Client: xxxarr.NewRadarrClient(url, apiKey, r)},
		application: "radarr",
		metrics:     createMetrics("radarr", url),
		transport:   r,
		logger:      slog.Default().With(slog.String("collector", "radarr")),
	}
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string) *Collector {
	r := httpclient.NewRoundTripper(
		httpclient.WithInstrumentedCache(sonarrCacheTable, cacheExpiry, cleanupInterval, "mediamon", "", "sonarr"),
		httpclient.WithMetrics("mediamon", "", "sonarr"),
	)

	return &Collector{
		Scraper:     scraper.SonarrScraper{Client: xxxarr.NewSonarrClient(url, apiKey, r)},
		application: "sonarr",
		metrics:     createMetrics("sonarr", url),
		transport:   r,
		logger:      slog.Default().With(slog.String("collector", "sonarr")),
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
	c.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	// TODO: http response's body.Close() sometimes panics in mediaclient/xxxarr ???
	defer func() {
		if err := recover(); err != nil {
			c.logger.Warn("scrape panicked", "err", err)
		}
	}()

	start := time.Now()
	stats, err := c.Scraper.Scrape(context.Background())
	if err != nil {
		// ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error", "Error getting "+c.application+" metrics", nil, nil), err)
		var err2 *xxxarr.ErrInvalidJSON
		if errors.As(err, &err2) {
			c.logger.Error("server returned invalid output", "err", err, "body", string(err2.Body))
		} else {
			c.logger.Error("failed to collect metrics", "err", err)
		}
		return
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["version"], prometheus.GaugeValue, float64(1), stats.Version)
	for key, value := range stats.Health {
		ch <- prometheus.MustNewConstMetric(c.metrics["health"], prometheus.GaugeValue, float64(value), key)
	}
	for _, title := range stats.Calendar {
		ch <- prometheus.MustNewConstMetric(c.metrics["calendar"], prometheus.GaugeValue, 1.0, title)
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["queued_count"], prometheus.GaugeValue, float64(len(stats.Queued)))

	totalBytes := make(map[string]float64)
	downloadedBytes := make(map[string]float64)
	for _, queued := range stats.Queued {
		totalBytes[queued.Name] += queued.TotalBytes
		downloadedBytes[queued.Name] += queued.DownloadedBytes
	}
	for name := range totalBytes {
		ch <- prometheus.MustNewConstMetric(c.metrics["queued_total"], prometheus.GaugeValue, totalBytes[name], name)
		ch <- prometheus.MustNewConstMetric(c.metrics["queued_downloaded"], prometheus.GaugeValue, downloadedBytes[name], name)
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["monitored"], prometheus.GaugeValue, float64(stats.Monitored))
	ch <- prometheus.MustNewConstMetric(c.metrics["unmonitored"], prometheus.GaugeValue, float64(stats.Unmonitored))
	c.transport.Collect(ch)

	c.logger.Debug("stats collected", "duration", time.Since(start))
}
