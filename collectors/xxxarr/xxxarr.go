package xxxarr

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
		prometheus.BuildFQName("mediamon", "xxxarr", "version"),
		"version info",
		[]string{"version", "url", "server"},
		nil,
	)

	calendarMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "calendar_count"),
		"Number of upcoming episodes / movies",
		[]string{"url", "server"},
		nil,
	)

	queuedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
		"Number of episodes / movies being downloaded",
		[]string{"url", "server"},
		nil,
	)

	monitoredMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
		"Number of monitored series / movies",
		[]string{"url", "server"},
		nil,
	)

	unmonitoredMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
		"Number of unmonitored series / movies",
		[]string{"url", "server"},
		nil,
	)
)

type Collector struct {
	mediaclient.XXXArrAPI
	cache.Cache
	application string
	url         string
}

type xxxArrStats struct {
	version     string
	calendar    int
	queued      int
	monitored   int
	unmonitored int
}

func NewCollector(url, apiKey, application string, interval time.Duration) prometheus.Collector {
	if application != "sonarr" && application != "radarr" {
		panic("invalid application: " + application)
	}

	c := &Collector{
		XXXArrAPI: &mediaclient.XXXArrClient{
			Client:      &http.Client{},
			URL:         url,
			APIKey:      apiKey,
			Application: application,
			Options: mediaclient.XXXArrOpts{
				PrometheusSummary: metrics.RequestDuration,
			},
		},
		application: application,
		url:         url,
	}

	c.Cache = *cache.New(interval, xxxArrStats{}, c.getStats)

	return c
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- calendarMetric
	ch <- queuedMetric
	ch <- monitoredMetric
	ch <- unmonitoredMetric
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(xxxArrStats)

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.url, coll.application)
	ch <- prometheus.MustNewConstMetric(calendarMetric, prometheus.GaugeValue, float64(stats.calendar), coll.url, coll.application)
	ch <- prometheus.MustNewConstMetric(queuedMetric, prometheus.GaugeValue, float64(stats.queued), coll.url, coll.application)
	ch <- prometheus.MustNewConstMetric(monitoredMetric, prometheus.GaugeValue, float64(stats.monitored), coll.url, coll.application)
	ch <- prometheus.MustNewConstMetric(unmonitoredMetric, prometheus.GaugeValue, float64(stats.unmonitored), coll.url, coll.application)
}

func (coll *Collector) getStats() (interface{}, error) {
	var stats xxxArrStats
	var err error

	ctx := context.Background()

	stats.version, err = coll.GetVersion(ctx)

	if err == nil {
		stats.calendar, err = coll.GetCalendar(ctx)
	}

	if err == nil {
		stats.queued, err = coll.GetQueue(ctx)
	}

	if err == nil {
		stats.monitored, stats.unmonitored, err = coll.GetMonitored(ctx)
	}

	return stats, err
}
