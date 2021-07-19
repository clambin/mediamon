package xxxarr

import (
	"context"
	"github.com/clambin/mediamon/cache"
	"github.com/clambin/mediamon/pkg/mediaclient"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type Collector struct {
	mediaclient.XXXArrAPI
	cache.Cache
	version     *prometheus.Desc
	calendar    *prometheus.Desc
	queued      *prometheus.Desc
	monitored   *prometheus.Desc
	unmonitored *prometheus.Desc
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
		},
		version: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "version"),
			"version info",
			[]string{"version"},
			prometheus.Labels{"url": url, "server": application},
		),
		calendar: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "calendar_count"),
			"Number of upcoming episodes / movies",
			nil,
			prometheus.Labels{"url": url, "server": application},
		),
		queued: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
			"Number of episodes / movies being downloaded",
			nil,
			prometheus.Labels{"url": url, "server": application},
		),
		monitored: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
			"Number of monitored series / movies",
			nil,
			prometheus.Labels{"url": url, "server": application},
		),
		unmonitored: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
			"Number of unmonitored series / movies",
			nil,
			prometheus.Labels{"url": url, "server": application},
		),
	}

	c.Cache = *cache.New(interval, xxxArrStats{}, c.getStats)

	return c
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- coll.version
	ch <- coll.calendar
	ch <- coll.queued
	ch <- coll.monitored
	ch <- coll.unmonitored
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(xxxArrStats)

	ch <- prometheus.MustNewConstMetric(coll.version, prometheus.GaugeValue, float64(1), stats.version)
	ch <- prometheus.MustNewConstMetric(coll.calendar, prometheus.GaugeValue, float64(stats.calendar))
	ch <- prometheus.MustNewConstMetric(coll.queued, prometheus.GaugeValue, float64(stats.queued))
	ch <- prometheus.MustNewConstMetric(coll.monitored, prometheus.GaugeValue, float64(stats.monitored))
	ch <- prometheus.MustNewConstMetric(coll.unmonitored, prometheus.GaugeValue, float64(stats.unmonitored))
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
