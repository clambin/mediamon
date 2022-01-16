package transmission

import (
	"context"
	metrics2 "github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/transmission"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	url string
}

type transmissionStats struct {
	version  string
	active   int
	paused   int
	download int
	upload   int
}

// NewCollector creates a new Collector
func NewCollector(url string, _ time.Duration) prometheus.Collector {
	return &Collector{
		API: &transmission.Client{
			Client: &http.Client{},
			URL:    url,
			Options: transmission.Options{
				PrometheusMetrics: metrics2.APIClientMetrics{
					Latency: metrics.Latency,
					Errors:  metrics.Errors,
				},
			},
		},
		url: url,
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := coll.getStats()
	if err != nil {
		log.WithError(err).Warning("failed to collect transmission metrics")
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.url)
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.active), coll.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.paused), coll.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.download), coll.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.upload), coll.url)
}

func (coll *Collector) getStats() (stats transmissionStats, err error) {
	ctx := context.Background()

	stats.version, err = coll.GetVersion(ctx)

	if err == nil {
		stats.active, stats.paused, stats.download, stats.upload, err = coll.GetStats(ctx)
	}

	return stats, err
}
