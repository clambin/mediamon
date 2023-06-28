package transmission

import (
	"context"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
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

// API interface
//
//go:generate mockery --name API
type API interface {
	GetSessionParameters(ctx context.Context) (transmission.SessionParameters, error)
	GetSessionStatistics(ctx context.Context) (stats transmission.SessionStats, err error)
}

// Collector presents Transmission statistics as Prometheus metrics
type Collector struct {
	API
	url       string
	transport *httpclient.RoundTripper
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
		API:       transmission.NewClient(url, r),
		url:       url,
		transport: r,
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
	coll.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	stats, err := coll.getStats()
	if err != nil {
		//ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error","Error getting transmission metrics", nil, nil),err)
		slog.Error("failed to collect transmission metrics", "err", err)
		return
	}
	stats.collect(ch, coll.url)
	coll.transport.Collect(ch)
	defer slog.Debug("transmission stats collected", "duration", time.Since(start))
}

func (coll *Collector) getStats() (stats transmissionStats, err error) {
	ctx := context.Background()
	if stats.version, err = coll.getVersion(ctx); err == nil {
		stats.active, stats.paused, stats.download, stats.upload, err = coll.getSessionStats(ctx)
	}
	return stats, err
}

func (coll *Collector) getVersion(ctx context.Context) (string, error) {
	params, err := coll.API.GetSessionParameters(ctx)
	return params.Arguments.Version, err
}

func (coll *Collector) getSessionStats(ctx context.Context) (int, int, int, int, error) {
	stats, err := coll.API.GetSessionStatistics(ctx)
	return stats.Arguments.ActiveTorrentCount, stats.Arguments.PausedTorrentCount, stats.Arguments.DownloadSpeed, stats.Arguments.UploadSpeed, err
}
