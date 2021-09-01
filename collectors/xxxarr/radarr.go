package xxxarr

import (
	"github.com/clambin/mediamon/cache"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient"
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

type RadarrCollector struct {
	Updater
	cache.Cache
}

func NewRadarrCollector(url, apiKey string, interval time.Duration) prometheus.Collector {
	collector := &RadarrCollector{
		Updater: Updater{
			XXXArrAPI: &mediaclient.XXXArrClient{
				Client:      &http.Client{},
				URL:         url,
				APIKey:      apiKey,
				Application: "radarr",
				Options: mediaclient.XXXArrOpts{
					PrometheusSummary: metrics.RequestDuration,
				},
			},
		},
	}
	collector.Cache = cache.Cache{
		Expiry:    interval,
		LastStats: xxxArrStats{},
		Updater:   collector.getStats,
	}
	return collector
}

func (coll *RadarrCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- radarrVersionMetric
	ch <- radarrCalendarMetric
	ch <- radarrQueuedMetric
	ch <- radarrMonitoredMetric
	ch <- radarrUnmonitoredMetric
}

func (coll *RadarrCollector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(xxxArrStats)

	ch <- prometheus.MustNewConstMetric(sonarrVersionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.Updater.XXXArrAPI.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrCalendarMetric, prometheus.GaugeValue, float64(stats.calendar), coll.Updater.XXXArrAPI.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrQueuedMetric, prometheus.GaugeValue, float64(stats.queued), coll.Updater.XXXArrAPI.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrMonitoredMetric, prometheus.GaugeValue, float64(stats.monitored), coll.Updater.XXXArrAPI.GetURL())
	ch <- prometheus.MustNewConstMetric(sonarrUnmonitoredMetric, prometheus.GaugeValue, float64(stats.unmonitored), coll.Updater.XXXArrAPI.GetURL())
}
