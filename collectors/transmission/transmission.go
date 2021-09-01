package transmission

import (
	"context"
	"github.com/clambin/mediamon/cache"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient"
	"github.com/prometheus/client_golang/prometheus"
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

type Collector struct {
	mediaclient.TransmissionAPI
	cache.Cache
	url string
}

type transmissionStats struct {
	version  string
	active   int
	paused   int
	download int
	upload   int
}

func NewCollector(url string, interval time.Duration) prometheus.Collector {
	c := &Collector{
		TransmissionAPI: &mediaclient.TransmissionClient{
			Client: &http.Client{},
			URL:    url,
			Options: mediaclient.TransmissionOpts{
				PrometheusSummary: metrics.RequestDuration,
			},
		},
		url: url,
	}
	c.Cache = *cache.New(interval, transmissionStats{}, c.getStats)
	return c
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(transmissionStats)

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.url)
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.active), coll.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.paused), coll.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.download), coll.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.upload), coll.url)
}

func (coll *Collector) getStats() (interface{}, error) {
	var stats transmissionStats
	var err error

	ctx := context.Background()

	stats.version, err = coll.GetVersion(ctx)

	if err == nil {
		stats.active, stats.paused, stats.download, stats.upload, err = coll.GetStats(ctx)
	}

	return stats, err
}
