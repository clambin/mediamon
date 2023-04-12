package transmission

import (
	"context"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/pkg/mediaclient/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
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

// Collector presents Transmission statistics as Prometheus metrics
type Collector struct {
	transmission.API
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

// NewCollector creates a new Collector
func NewCollector(url string) *Collector {
	r := httpclient.NewRoundTripper(httpclient.WithRoundTripperMetrics{Namespace: "mediamon", Application: "transmission"})
	return &Collector{
		API: &transmission.Client{
			HTTPClient: &http.Client{Transport: r},
			URL:        url,
		},
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
		/*
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("mediamon_error",
					"Error getting transmission metrics", nil, nil),
				err)
		*/
		slog.Error("failed to collect transmission metrics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.url)
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.active), coll.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.paused), coll.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.download), coll.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.upload), coll.url)
	coll.transport.Collect(ch)
	defer slog.Debug("transmission stats collected", "duration", time.Since(start))
}

func (coll *Collector) getStats() (stats transmissionStats, err error) {
	ctx := context.Background()

	stats.version, err = coll.getVersion(ctx)

	if err == nil {
		stats.active, stats.paused, stats.download, stats.upload, err = coll.getSessionStats(ctx)
	}

	return stats, err
}

func (coll *Collector) getVersion(ctx context.Context) (version string, err error) {
	var params transmission.SessionParameters
	params, err = coll.API.GetSessionParameters(ctx)
	if err == nil {
		version = params.Arguments.Version
	}
	return
}

func (coll *Collector) getSessionStats(ctx context.Context) (active int, paused int, download int, upload int, err error) {
	var stats transmission.SessionStats
	stats, err = coll.API.GetSessionStatistics(ctx)
	if err == nil {
		active = stats.Arguments.ActiveTorrentCount
		paused = stats.Arguments.PausedTorrentCount
		download = stats.Arguments.DownloadSpeed
		upload = stats.Arguments.UploadSpeed
	}
	return
}
