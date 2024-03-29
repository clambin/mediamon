package transmission

import (
	"context"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"sync"
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

var _ prometheus.Collector = &Collector{}

// NewCollector creates a new Collector
func NewCollector(url string) *Collector {
	m := metrics.NewRequestSummaryMetrics("mediamon", "", map[string]string{"application": "transmission"})
	r := roundtripper.New(roundtripper.WithRequestMetrics(m))
	return &Collector{
		Transmission: transmission.NewClient(url, r),
		url:          url,
		metrics:      m,
		logger:       slog.Default().With(slog.String("collector", "transmission")),
	}
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

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); c.collectVersion(ch) }()
	go func() { defer wg.Done(); c.collectStats(ch) }()
	wg.Wait()
	c.metrics.Collect(ch)
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) {
	params, err := c.Transmission.GetSessionParameters(context.Background())
	if err != nil {
		c.logger.Error("failed to get session parameters", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), params.Arguments.Version, c.url)
}

func (c *Collector) collectStats(ch chan<- prometheus.Metric) {
	stats, err := c.Transmission.GetSessionStatistics(context.Background())
	if err != nil {
		c.logger.Error("failed to get session statistics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.Arguments.ActiveTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.Arguments.PausedTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.Arguments.DownloadSpeed), c.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.Arguments.UploadSpeed), c.url)
}
