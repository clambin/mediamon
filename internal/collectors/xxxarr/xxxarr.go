package xxxarr

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type QueuedItem struct {
	Name            string
	TotalBytes      int64
	DownloadedBytes int64
}

type Library struct {
	Monitored   int
	Unmonitored int
}

func WithToken(token string) func(ctx context.Context, req *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		if token == "" {
			return fmt.Errorf("no token provided")
		}
		req.Header.Set("X-Api-Key", token)
		return nil
	}
}

type Collector struct {
	client      Client
	metrics     map[string]*prometheus.Desc
	logger      *slog.Logger
	application string
}

// Client presents a unified interface to Sonarr/Radarr clients
type Client interface {
	GetVersion(context.Context) (string, error)
	GetHealth(context.Context) (map[string]int, error)
	GetCalendar(context.Context, int) ([]string, error)
	GetQueue(context.Context) ([]QueuedItem, error)
	GetLibrary(context.Context) (Library, error)
}

var (
	_ Client = Sonarr{}
	_ Client = Radarr{}
)

func NewRadarrCollector(url, apiKey string, httpClient *http.Client, logger *slog.Logger) (prometheus.Collector, error) {
	client, err := NewRadarrClient(url, apiKey, httpClient)
	if err != nil {
		return nil, fmt.Errorf("create radarr client: %w", err)
	}

	c := Collector{
		client:      client,
		application: "radarr",
		metrics:     createMetrics("radarr", url),
		logger:      logger,
	}
	return &c, nil
}

func NewSonarrCollector(url, apiKey string, httpClient *http.Client, logger *slog.Logger) (prometheus.Collector, error) {
	client, err := NewSonarrClient(url, apiKey, httpClient)
	if err != nil {
		return nil, fmt.Errorf("create sonarr client: %w", err)
	}

	c := Collector{
		client:      client,
		application: "sonarr",
		metrics:     createMetrics("sonarr", url),
		logger:      logger,
	}
	return &c, nil
}

// Describe implements the prometheus.Collector interface
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
}

// Collect implements the prometheus.Collector interface
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	var g sync.WaitGroup
	g.Go(func() { c.collectVersion(ch) })
	g.Go(func() { c.collectHealth(ch) })
	g.Go(func() { c.collectCalendar(ch) })
	g.Go(func() { c.collectQueue(ch) })
	g.Go(func() { c.collectLibrary(ch) })
	g.Wait()
}

func (c Collector) collectVersion(ch chan<- prometheus.Metric) {
	version, err := c.client.GetVersion(context.Background())
	if err != nil {
		c.logger.Error("failed to get version", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.metrics["version"], prometheus.GaugeValue, float64(1), version)
}

func (c Collector) collectHealth(ch chan<- prometheus.Metric) {
	health, err := c.client.GetHealth(context.Background())
	if err != nil {
		c.logger.Error("failed to get health", "err", err)
		return
	}
	for key, value := range health {
		ch <- prometheus.MustNewConstMetric(c.metrics["health"], prometheus.GaugeValue, float64(value), key)
	}
}

func (c Collector) collectCalendar(ch chan<- prometheus.Metric) {
	calendar, err := c.client.GetCalendar(context.Background(), 1)
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

func (c Collector) collectQueue(ch chan<- prometheus.Metric) {
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

func (c Collector) collectLibrary(ch chan<- prometheus.Metric) {
	library, err := c.client.GetLibrary(context.Background())
	if err != nil {
		c.logger.Error("failed to get library", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.metrics["monitored"], prometheus.GaugeValue, float64(library.Monitored))
	ch <- prometheus.MustNewConstMetric(c.metrics["unmonitored"], prometheus.GaugeValue, float64(library.Unmonitored))
}

func createMetrics(application, url string) map[string]*prometheus.Desc {
	constLabels := prometheus.Labels{
		"application": application,
		"url":         url,
	}
	return map[string]*prometheus.Desc{
		"version": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "version"),
			"Version info",
			[]string{"version"},
			constLabels,
		),
		"health": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "health"),
			"Server health",
			[]string{"type"},
			constLabels,
		),
		"calendar": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "calendar"),
			"Upcoming episodes / movies",
			[]string{"title"},
			constLabels,
		),
		"queued_count": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
			"Episodes / movies being downloaded",
			nil,
			constLabels,
		),
		"queued_total": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_total_bytes"),
			"Size of episode / movie being downloaded in bytes",
			[]string{"title"},
			constLabels,
		),
		"queued_downloaded": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_downloaded_bytes"),
			"Downloaded size of episode / movie being downloaded in bytes",
			[]string{"title"},
			constLabels,
		),
		"monitored": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
			"Number of Monitored series / movies",
			nil,
			constLabels,
		),
		"unmonitored": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
			"Number of Unmonitored series / movies",
			nil,
			constLabels,
		),
	}
}
