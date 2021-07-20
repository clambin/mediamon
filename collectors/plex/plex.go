package plex

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
	mediaclient.PlexAPI
	cache.Cache
	version     *prometheus.Desc
	users       *prometheus.Desc
	modes       *prometheus.Desc
	transcoding *prometheus.Desc
	speed       *prometheus.Desc
}

type plexStats struct {
	version     string
	users       map[string]int
	modes       map[string]int
	transcoding int
	speed       float64
}

func NewCollector(url, username, password string, interval time.Duration) prometheus.Collector {
	c := &Collector{
		PlexAPI: &mediaclient.PlexClient{
			Client:   &http.Client{},
			URL:      url,
			UserName: username,
			Password: password,
			Options: mediaclient.PlexOpts{
				PrometheusSummary: metrics.RequestDuration,
			},
		},
		version: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "plex", "version"),
			"version info",
			[]string{"version"},
			prometheus.Labels{"url": url},
		),
		transcoding: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "plex", "transcoder_encoding_count"),
			"Number of active transcoders",
			nil,
			prometheus.Labels{"url": url},
		),
		speed: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "plex", "transcoder_speed_total"),
			"Speed of active transcoders",
			nil,
			prometheus.Labels{"url": url},
		),
		users: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "plex", "session_count"),
			"Active Plex Sessions",
			[]string{"user"},
			prometheus.Labels{"url": url},
		),
		modes: prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "plex", "transcoder_type_count"),
			"Active Transcoder count by type",
			[]string{"mode"},
			prometheus.Labels{"url": url},
		),
	}

	c.Cache = *cache.New(interval, plexStats{}, c.getStats)

	return c
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- coll.version
	ch <- coll.transcoding
	ch <- coll.speed
	ch <- coll.users
	ch <- coll.modes
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(plexStats)

	ch <- prometheus.MustNewConstMetric(coll.version, prometheus.GaugeValue, float64(1), stats.version)
	ch <- prometheus.MustNewConstMetric(coll.transcoding, prometheus.GaugeValue, float64(stats.transcoding))
	ch <- prometheus.MustNewConstMetric(coll.speed, prometheus.GaugeValue, stats.speed)
	for user, count := range stats.users {
		ch <- prometheus.MustNewConstMetric(coll.users, prometheus.GaugeValue, float64(count), user)
	}
	for mode, count := range stats.modes {
		ch <- prometheus.MustNewConstMetric(coll.modes, prometheus.GaugeValue, float64(count), mode)
	}
}

func (coll *Collector) getStats() (interface{}, error) {
	var stats plexStats
	var err error

	ctx := context.Background()

	stats.version, err = coll.GetVersion(ctx)

	if err == nil {
		stats.users, stats.modes, stats.transcoding, stats.speed, err = coll.GetSessions(ctx)
	}

	return stats, err
}
