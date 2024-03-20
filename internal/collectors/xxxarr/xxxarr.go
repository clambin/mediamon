package xxxarr

import (
	"context"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/xxxarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients"
	"github.com/clambin/mediamon/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"sync"
	"time"
)

// Collector presents Sonarr/Radarr statistics as Prometheus metrics
type Collector struct {
	client       Client
	application  string
	metrics      map[string]*prometheus.Desc
	tpMetrics    roundtripper.RoundTripMetrics
	cacheMetrics roundtripper.CacheMetrics
	logger       *slog.Logger
}

type Client interface {
	GetVersion(context.Context) (string, error)
	GetHealth(context.Context) (map[string]int, error)
	GetCalendar(context.Context) ([]string, error)
	GetQueue(context.Context) ([]clients.QueuedItem, error)
	GetLibrary(context.Context) (clients.Library, error)
}

var (
	radarrCacheTable = roundtripper.CacheTable{
		{Path: `/api/v3/system/status`, Expiry: time.Minute},
		{Path: `/api/v3/calendar`, Expiry: time.Minute},
		{Path: `/api/v3/movie`},
		{Path: `/api/v3/movie/[\d+]`, IsRegExp: true},
	}

	sonarrCacheTable = roundtripper.CacheTable{
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
	tpMetrics := metrics.NewCustomizedRoundTripMetrics("mediamon", "", "radarr", chopPath)
	cacheMetrics := metrics.NewCustomizedCacheMetrics("mediamon", "", "radarr", chopPath)

	r := roundtripper.New(
		roundtripper.WithInstrumentedCache(radarrCacheTable, cacheExpiry, cleanupInterval, cacheMetrics),
		roundtripper.WithInstrumentedRoundTripper(tpMetrics),
	)

	return &Collector{
		client:       clients.Radarr{Client: xxxarr.NewRadarrClient(url, apiKey, r)},
		application:  "radarr",
		metrics:      createMetrics("radarr", url),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       slog.Default().With(slog.String("collector", "radarr")),
	}
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string) *Collector {
	tpMetrics := metrics.NewCustomizedRoundTripMetrics("mediamon", "", "sonarr", chopPath)
	cacheMetrics := metrics.NewCustomizedCacheMetrics("mediamon", "", "sonarr", chopPath)

	r := roundtripper.New(
		roundtripper.WithInstrumentedCache(sonarrCacheTable, cacheExpiry, cleanupInterval, cacheMetrics),
		roundtripper.WithInstrumentedRoundTripper(tpMetrics),
	)

	return &Collector{
		client:       clients.Sonarr{Client: xxxarr.NewSonarrClient(url, apiKey, r)},
		application:  "sonarr",
		metrics:      createMetrics("sonarr", url),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       slog.Default().With(slog.String("collector", "sonarr")),
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
	c.tpMetrics.Describe(ch)
	c.cacheMetrics.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(5)
	go func() { defer wg.Done(); c.collectVersion(ch) }()
	go func() { defer wg.Done(); c.collectHealth(ch) }()
	go func() { defer wg.Done(); c.collectCalendar(ch) }()
	go func() { defer wg.Done(); c.collectQueue(ch) }()
	go func() { defer wg.Done(); c.collectLibrary(ch) }()
	wg.Wait()
	c.tpMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) {
	version, err := c.client.GetVersion(context.Background())
	if err != nil {
		c.logger.Error("failed to get version", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.metrics["version"], prometheus.GaugeValue, float64(1), version)
}

func (c *Collector) collectHealth(ch chan<- prometheus.Metric) {
	health, err := c.client.GetHealth(context.Background())
	if err != nil {
		c.logger.Error("failed to get health", "err", err)
		return
	}
	for key, value := range health {
		ch <- prometheus.MustNewConstMetric(c.metrics["health"], prometheus.GaugeValue, float64(value), key)
	}
}

func (c *Collector) collectCalendar(ch chan<- prometheus.Metric) {
	calendar, err := c.client.GetCalendar(context.Background())
	if err != nil {
		c.logger.Error("failed to get calendar", "err", err)
		return
	}
	for name, count := range groupNames(calendar) {
		ch <- prometheus.MustNewConstMetric(c.metrics["calendar"], prometheus.GaugeValue, float64(count), name)
	}
}

func groupNames(names []string) map[string]int {
	result := make(map[string]int)
	for i := range names {
		result[names[i]]++
	}
	return result
}

func (c *Collector) collectQueue(ch chan<- prometheus.Metric) {
	queue, err := c.client.GetQueue(context.Background())
	if err != nil {
		c.logger.Error("failed to get queue", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["queued_count"], prometheus.GaugeValue, float64(len(queue)))

	totalBytes := make(map[string]int64)
	downloadedBytes := make(map[string]int64)
	for _, queued := range queue {
		totalBytes[queued.Name] += queued.TotalBytes
		downloadedBytes[queued.Name] += queued.DownloadedBytes
	}
	for name := range totalBytes {
		ch <- prometheus.MustNewConstMetric(c.metrics["queued_total"], prometheus.GaugeValue, float64(totalBytes[name]), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["queued_downloaded"], prometheus.GaugeValue, float64(downloadedBytes[name]), name)
	}
}

func (c *Collector) collectLibrary(ch chan<- prometheus.Metric) {
	library, err := c.client.GetLibrary(context.Background())
	if err != nil {
		c.logger.Error("failed to get library", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.metrics["monitored"], prometheus.GaugeValue, float64(library.Monitored))
	ch <- prometheus.MustNewConstMetric(c.metrics["unmonitored"], prometheus.GaugeValue, float64(library.Unmonitored))
}
