package xxxarr

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/xxxarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients"
	collectorBreaker "github.com/clambin/mediamon/v2/pkg/collector-breaker"
	customMetrics "github.com/clambin/mediamon/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

// Collector presents Sonarr/Radarr statistics as Prometheus metrics
type Collector struct {
	client       Client
	application  string
	metrics      map[string]*prometheus.Desc
	tpMetrics    metrics.RequestMetrics
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
func NewRadarrCollector(url, apiKey string, logger *slog.Logger) *collectorBreaker.CBCollector {
	tpMetrics := customMetrics.NewCustomizedRoundTripMetrics("mediamon", "", map[string]string{"application": "radarr"}, chopPath)
	cacheMetrics := customMetrics.NewCustomizedCacheMetrics("mediamon", "", "radarr", chopPath)

	r := roundtripper.New(
		roundtripper.WithInstrumentedCache(radarrCacheTable, cacheExpiry, cleanupInterval, cacheMetrics),
		roundtripper.WithRequestMetrics(tpMetrics),
	)

	c := Collector{
		client:       clients.Radarr{Client: xxxarr.NewRadarrClient(url, apiKey, r)},
		application:  "radarr",
		metrics:      createMetrics("radarr", url),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       logger,
	}
	return collectorBreaker.New("radarr", &c, logger)
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string, logger *slog.Logger) *collectorBreaker.CBCollector {
	tpMetrics := customMetrics.NewCustomizedRoundTripMetrics("mediamon", "", map[string]string{"application": "sonarr"}, chopPath)
	cacheMetrics := customMetrics.NewCustomizedCacheMetrics("mediamon", "", "sonarr", chopPath)

	r := roundtripper.New(
		roundtripper.WithInstrumentedCache(sonarrCacheTable, cacheExpiry, cleanupInterval, cacheMetrics),
		roundtripper.WithRequestMetrics(tpMetrics),
	)

	c := Collector{
		client:       clients.Sonarr{Client: xxxarr.NewSonarrClient(url, apiKey, r)},
		application:  "sonarr",
		metrics:      createMetrics("sonarr", url),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       logger,
	}
	return collectorBreaker.New("sonarr", &c, logger)
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
	c.tpMetrics.Describe(ch)
	c.cacheMetrics.Describe(ch)
}

// CollectE implements the prometheus.Collector interface
func (c *Collector) CollectE(ch chan<- prometheus.Metric) error {
	var g errgroup.Group
	g.Go(func() error { return c.collectVersion(ch) })
	g.Go(func() error { return c.collectHealth(ch) })
	g.Go(func() error { return c.collectCalendar(ch) })
	g.Go(func() error { return c.collectQueue(ch) })
	g.Go(func() error { return c.collectLibrary(ch) })
	err := g.Wait()
	c.tpMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
	return err
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) error {
	version, err := c.client.GetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(c.metrics["version"], prometheus.GaugeValue, float64(1), version)
	return nil
}

func (c *Collector) collectHealth(ch chan<- prometheus.Metric) error {
	health, err := c.client.GetHealth(context.Background())
	if err != nil {
		return fmt.Errorf("health: %w", err)
	}
	for key, value := range health {
		ch <- prometheus.MustNewConstMetric(c.metrics["health"], prometheus.GaugeValue, float64(value), key)
	}
	return nil
}

func (c *Collector) collectCalendar(ch chan<- prometheus.Metric) error {
	calendar, err := c.client.GetCalendar(context.Background())
	if err != nil {
		return fmt.Errorf("calendar: %w", err)
	}
	for name, count := range groupNames(calendar) {
		ch <- prometheus.MustNewConstMetric(c.metrics["calendar"], prometheus.GaugeValue, float64(count), name)
	}
	return nil
}

func groupNames(names []string) map[string]int {
	result := make(map[string]int)
	for i := range names {
		result[names[i]]++
	}
	return result
}

func (c *Collector) collectQueue(ch chan<- prometheus.Metric) error {
	queue, err := c.client.GetQueue(context.Background())
	if err != nil {
		return fmt.Errorf("queue: %w", err)
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
	return nil
}

func (c *Collector) collectLibrary(ch chan<- prometheus.Metric) error {
	library, err := c.client.GetLibrary(context.Background())
	if err != nil {
		return fmt.Errorf("library: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(c.metrics["monitored"], prometheus.GaugeValue, float64(library.Monitored))
	ch <- prometheus.MustNewConstMetric(c.metrics["unmonitored"], prometheus.GaugeValue, float64(library.Unmonitored))
	return nil
}
