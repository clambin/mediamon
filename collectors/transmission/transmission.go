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

type Collector struct {
	mediaclient.TransmissionAPI
	cache.Cache
	version        *prometheus.Desc
	activeTorrents *prometheus.Desc
	pausedTorrents *prometheus.Desc
	downloadSpeed  *prometheus.Desc
	uploadSpeed    *prometheus.Desc
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
		version: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "transmission", "version"),
			"version info",
			[]string{"version"},
			prometheus.Labels{"url": url},
		),
		activeTorrents: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "transmission", "active_torrent_count"),
			"Number of active torrents",
			nil,
			prometheus.Labels{"url": url},
		),
		pausedTorrents: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "transmission", "paused_torrent_count"),
			"Number of paused torrents",
			nil,
			prometheus.Labels{"url": url},
		),
		downloadSpeed: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "transmission", "download_speed"),
			"Transmission download speed in bytes / sec",
			nil,
			prometheus.Labels{"url": url},
		),
		uploadSpeed: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "transmission", "upload_speed"),
			"Transmission upload speed in bytes / sec",
			nil,
			prometheus.Labels{"url": url},
		),
	}
	c.Cache = *cache.New(interval, transmissionStats{}, c.getStats)
	return c
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- coll.version
	ch <- coll.activeTorrents
	ch <- coll.pausedTorrents
	ch <- coll.downloadSpeed
	ch <- coll.uploadSpeed
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(transmissionStats)

	ch <- prometheus.MustNewConstMetric(coll.version, prometheus.GaugeValue, float64(1), stats.version)
	ch <- prometheus.MustNewConstMetric(coll.activeTorrents, prometheus.GaugeValue, float64(stats.active))
	ch <- prometheus.MustNewConstMetric(coll.pausedTorrents, prometheus.GaugeValue, float64(stats.paused))
	ch <- prometheus.MustNewConstMetric(coll.downloadSpeed, prometheus.GaugeValue, float64(stats.download))
	ch <- prometheus.MustNewConstMetric(coll.uploadSpeed, prometheus.GaugeValue, float64(stats.upload))
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
