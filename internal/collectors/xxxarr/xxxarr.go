package xxxarr

import (
	"codeberg.org/clambin/go-common/httputils/metrics"
	"codeberg.org/clambin/go-common/httputils/roundtripper"
	"context"
	"fmt"
	collectorbreaker "github.com/clambin/mediamon/v2/collector-breaker"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"strconv"
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
	GetCalendar(context.Context, int) ([]string, error)
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
		{Path: `/api/v3/episode`, IsRegExp: true},
	}
)

const (
	cacheExpiry     = 15 * time.Minute
	cleanupInterval = 5 * time.Minute
)

// NewRadarrCollector creates a new RadarrCollector
func NewRadarrCollector(url, apiKey string, logger *slog.Logger) (*collectorbreaker.CBCollector, error) {
	tpMetrics := metrics.NewRequestMetrics(metrics.Options{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "radarr"},
		LabelValues: func(request *http.Request, i int) (method string, path string, code string) {
			return request.Method, choppedPath(request), strconv.Itoa(i)
		},
	})
	cacheMetrics := roundtripper.NewCacheMetrics(roundtripper.CacheMetricsOptions{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "radarr"},
		GetPath:     choppedPath,
	})

	httpClient := http.Client{
		Transport: roundtripper.New(
			roundtripper.WithCache(roundtripper.CacheOptions{
				CacheTable:        radarrCacheTable,
				DefaultExpiration: cacheExpiry,
				CleanupInterval:   cleanupInterval,
				CacheMetrics:      cacheMetrics,
				GetKey: func(r *http.Request) string {
					return http.MethodGet + "|" + r.URL.Path + "|" + r.URL.Query().Encode()
				},
			}),
			roundtripper.WithRequestMetrics(tpMetrics),
		),
	}

	client, err := clients.NewRadarrClient(url, apiKey, &httpClient)
	if err != nil {
		return nil, fmt.Errorf("create radarr client: %w", err)
	}

	c := Collector{
		client:       client,
		application:  "radarr",
		metrics:      createMetrics("radarr", url),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       logger,
	}
	return collectorbreaker.New("radarr", &c, logger), nil
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(target, apiKey string, logger *slog.Logger) (*collectorbreaker.CBCollector, error) {
	tpMetrics := metrics.NewRequestMetrics(metrics.Options{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "sonarr"},
		LabelValues: func(request *http.Request, i int) (method string, path string, code string) {
			return request.Method, choppedPath(request), strconv.Itoa(i)
		},
	})
	cacheMetrics := roundtripper.NewCacheMetrics(roundtripper.CacheMetricsOptions{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "sonarr"},
		GetPath:     choppedPath,
	})

	httpClient := http.Client{
		Transport: roundtripper.New(
			roundtripper.WithCache(roundtripper.CacheOptions{
				CacheTable:        sonarrCacheTable,
				DefaultExpiration: cacheExpiry,
				CleanupInterval:   cleanupInterval,
				CacheMetrics:      cacheMetrics,
				GetKey: func(r *http.Request) string {
					return http.MethodGet + "|" + r.URL.Path + "|" + r.URL.Query().Encode()
				},
			}),
			roundtripper.WithRequestMetrics(tpMetrics),
		),
	}

	client, err := clients.NewSonarrClient(target, apiKey, &httpClient)
	if err != nil {
		return nil, fmt.Errorf("create sonarr client: %w", err)
	}

	c := Collector{
		client:       client,
		application:  "sonarr",
		metrics:      createMetrics("sonarr", target),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       logger,
	}
	return collectorbreaker.New("sonarr", &c, logger), nil
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
	calendar, err := c.client.GetCalendar(context.Background(), 1)
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
