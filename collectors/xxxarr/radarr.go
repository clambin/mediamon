package xxxarr

import (
	"github.com/clambin/mediamon/cache"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	metrics2 "github.com/clambin/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	radarrVersionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "version"),
		"version info",
		[]string{"version", "url"},
		prometheus.Labels{"application": "radarr"},
	)

	radarrCalendarMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "calendar_count"),
		"Number of upcoming episodes / movies",
		[]string{"url"},
		prometheus.Labels{"application": "radarr"},
	)

	radarrQueuedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
		"Number of episodes / movies being downloaded",
		[]string{"url"},
		prometheus.Labels{"application": "radarr"},
	)

	radarrMonitoredMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
		"Number of monitored series / movies",
		[]string{"url"},
		prometheus.Labels{"application": "radarr"},
	)

	radarrUnmonitoredMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
		"Number of unmonitored series / movies",
		[]string{"url"},
		prometheus.Labels{"application": "radarr"},
	)
)

// RadarrCollector presents Radarr statistics as Prometheus metrics
type RadarrCollector struct {
	Updater
	cache.Cache
}

// NewRadarrCollector creates a new RadarrCollector
func NewRadarrCollector(url, apiKey string, interval time.Duration) prometheus.Collector {
	collector := &RadarrCollector{
		Updater: Updater{
			API: &xxxarr.Client{
				Client:      &http.Client{},
				URL:         url,
				APIKey:      apiKey,
				Application: "radarr",
				Options: xxxarr.Options{
					PrometheusMetrics: metrics2.APIClientMetrics{
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
func (coll *RadarrCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- radarrVersionMetric
	ch <- radarrCalendarMetric
	ch <- radarrQueuedMetric
	ch <- radarrMonitoredMetric
	ch <- radarrUnmonitoredMetric
}

// Collect implements the prometheus.Collector interface
func (coll *RadarrCollector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(xxxArrStats)

	ch <- prometheus.MustNewConstMetric(radarrVersionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(radarrCalendarMetric, prometheus.GaugeValue, float64(stats.calendar), coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(radarrQueuedMetric, prometheus.GaugeValue, float64(stats.queued), coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(radarrMonitoredMetric, prometheus.GaugeValue, float64(stats.monitored), coll.Updater.API.GetURL())
	ch <- prometheus.MustNewConstMetric(radarrUnmonitoredMetric, prometheus.GaugeValue, float64(stats.unmonitored), coll.Updater.API.GetURL())
}
