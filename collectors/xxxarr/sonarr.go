package xxxarr

import (
	"github.com/clambin/mediamon/cache"
	"github.com/clambin/mediamon/metrics"
	metrics2 "github.com/clambin/mediamon/pkg/mediaclient/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	sonarrVersionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "version"),
		"version info",
		[]string{"version", "url"},
		prometheus.Labels{"application": "sonarr"},
	)

	sonarrCalendarMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "calendar_count"),
		"Number of upcoming episodes / movies",
		[]string{"url"},
		prometheus.Labels{"application": "sonarr"},
	)

	sonarrQueuedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
		"Number of episodes / movies being downloaded",
		[]string{"url"},
		prometheus.Labels{"application": "sonarr"},
	)

	sonarrMonitoredMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
		"Number of monitored series / movies",
		[]string{"url"},
		prometheus.Labels{"application": "sonarr"},
	)

	sonarrUnmonitoredMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
		"Number of unmonitored series / movies",
		[]string{"url"},
		prometheus.Labels{"application": "sonarr"},
	)
)

// SonarrCollector presents Sonarr statistics as Prometheus metrics
type SonarrCollector struct {
	Updater
	cache.Cache
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string, interval time.Duration) prometheus.Collector {
	collector := &SonarrCollector{
		Updater: Updater{
			API: &xxxarr.Client{
				Client:      &http.Client{},
				URL:         url,
				APIKey:      apiKey,
				Application: "sonarr",
				Options: xxxarr.Options{
					PrometheusMetrics: metrics2.PrometheusMetrics{
						Latency: metrics.Latency,
						Errors:  metrics.Errors,
					},
				},
			},
		},
	}
	collector.Cache = cache.Cache{
		Duration:  interval,
		LastStats: xxxArrStats{},
		Updater:   collector.getStats,
	}
	return collector
}

// Describe implements the prometheus.Collector interface
func (coll *SonarrCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sonarrVersionMetric
	ch <- sonarrCalendarMetric
	ch <- sonarrQueuedMetric
	ch <- sonarrMonitoredMetric
	ch <- sonarrUnmonitoredMetric
}

// Collect implements the prometheus.Collector interface
func (coll *SonarrCollector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(xxxArrStats)

	ch <- prometheus.MustNewConstMetric(sonarrVersionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrCalendarMetric, prometheus.GaugeValue, float64(stats.calendar), coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrQueuedMetric, prometheus.GaugeValue, float64(stats.queued), coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrMonitoredMetric, prometheus.GaugeValue, float64(stats.monitored), coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrUnmonitoredMetric, prometheus.GaugeValue, float64(stats.unmonitored), coll.Updater.API.GetURL())
}
