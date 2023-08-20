package transmission

import (
	"context"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediaclients/transmission"
	"github.com/prometheus/client_golang/prometheus"
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
	Getter
	url       string
	transport *httpclient.RoundTripper
	logger    *slog.Logger
}

var _ prometheus.Collector = &Collector{}

// Config items for Transmission collector
type Config struct {
	URL string
}

type transmissionStats struct {
	version  string
	active   int
	paused   int
	download int
	upload   int
}

func (s transmissionStats) collect(ch chan<- prometheus.Metric, url string) {
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), s.version, url)
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(s.active), url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(s.paused), url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(s.download), url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(s.upload), url)
}

// NewCollector creates a new Collector
func NewCollector(url string) *Collector {
	r := httpclient.NewRoundTripper(httpclient.WithMetrics("mediamon", "", "transmission"))
	return &Collector{
		Getter:    transmission.NewClient(url, r),
		url:       url,
		transport: r,
		logger:    slog.Default().With(slog.String("collector", "transmission")),
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
	c.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	stats, err := c.getStats()
	if err != nil {
		//ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error","Error getting transmission metrics", nil, nil),err)
		c.logger.Error("failed to collect stats", "err", err)
		return
	}
	stats.collect(ch, c.url)
	c.transport.Collect(ch)
	c.logger.Debug("stats collected", "duration", time.Since(start))
}

func (c *Collector) getStats() (stats transmissionStats, err error) {
	ctx := context.Background()
	if stats.version, err = c.getVersion(ctx); err == nil {
		stats.active, stats.paused, stats.download, stats.upload, err = c.getSessionStats(ctx)
	}
	return stats, err
}

func (c *Collector) getVersion(ctx context.Context) (string, error) {
	params, err := c.Getter.GetSessionParameters(ctx)
	return params.Arguments.Version, err
}

func (c *Collector) getSessionStats(ctx context.Context) (int, int, int, int, error) {
	stats, err := c.Getter.GetSessionStatistics(ctx)
	return stats.Arguments.ActiveTorrentCount, stats.Arguments.PausedTorrentCount, stats.Arguments.DownloadSpeed, stats.Arguments.UploadSpeed, err
}
