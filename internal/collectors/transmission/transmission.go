package transmission

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/transmission"
	"github.com/clambin/mediamon/v2/pkg/breaker"
	collectorBreaker "github.com/clambin/mediamon/v2/pkg/collector-breaker"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

var (
	versionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "version"),
		"version info",
		[]string{"version", "url"},
		nil,
	)

	activeTorrentsMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "active_torrent_count"),
		"Number of active torrents",
		[]string{"url"},
		nil,
	)

	pausedTorrentsMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "paused_torrent_count"),
		"Number of paused torrents",
		[]string{"url"},
		nil,
	)

	downloadSpeedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "download_speed"),
		"Transmission download speed in bytes / sec",
		[]string{"url"},
		nil,
	)

	uploadSpeedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "upload_speed"),
		"Transmission upload speed in bytes / sec",
		[]string{"url"},
		nil,
	)
)

// Getter interface
type Getter interface {
	GetSessionParameters(ctx context.Context) (transmission.SessionParameters, error)
	GetSessionStatistics(ctx context.Context) (stats transmission.SessionStats, err error)
}

// Collector presents Transmission statistics as Prometheus metrics
type Collector struct {
	Transmission Getter
	url          string
	metrics      metrics.RequestMetrics
	logger       *slog.Logger
}

var breakerConfiguration = breaker.Configuration{
	FailureThreshold: 2,
	OpenDuration:     time.Minute,
	SuccessThreshold: 2,
}

var _ collectorBreaker.Collector = &Collector{}

// NewCollector creates a new Collector
func NewCollector(url string, logger *slog.Logger) *collectorBreaker.CBCollector {
	m := metrics.NewRequestSummaryMetrics("mediamon", "", map[string]string{"application": "transmission"})
	c := Collector{
		Transmission: transmission.NewClient(url, roundtripper.New(roundtripper.WithRequestMetrics(m))),
		url:          url,
		metrics:      m,
		logger:       logger,
	}
	return collectorBreaker.New(&c, breakerConfiguration, logger)
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
	c.metrics.Describe(ch)
}

// CollectE implements the prometheus.Collector interface
func (c *Collector) CollectE(ch chan<- prometheus.Metric) error {
	var g errgroup.Group
	g.Go(func() error { return c.collectVersion(ch) })
	g.Go(func() error { return c.collectStats(ch) })
	err := g.Wait()
	c.metrics.Collect(ch)
	return err
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) error {
	params, err := c.Transmission.GetSessionParameters(context.Background())
	if err != nil {
		return fmt.Errorf("error getting session parameters: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), params.Arguments.Version, c.url)
	return nil
}

func (c *Collector) collectStats(ch chan<- prometheus.Metric) error {
	stats, err := c.Transmission.GetSessionStatistics(context.Background())
	if err != nil {
		return fmt.Errorf("error getting session statistics: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.Arguments.ActiveTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.Arguments.PausedTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.Arguments.DownloadSpeed), c.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.Arguments.UploadSpeed), c.url)
	return nil
}
