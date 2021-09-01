package plex

import (
	"context"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	versionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "version"),
		"version info",
		[]string{"version", "url"},
		nil,
	)

	transcodingMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_encoding_count"),
		"Number of active transcoders",
		[]string{"url"},
		nil,
	)

	speedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_speed_total"),
		"Speed of active transcoders",
		[]string{"url"},
		nil,
	)

	usersMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_count"),
		"Active Plex Sessions",
		[]string{"user", "url"},
		nil,
	)

	modesMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_type_count"),
		"Active Transcoder count by type",
		[]string{"mode", "url"},
		nil,
	)
)

type Collector struct {
	mediaclient.PlexAPI
	url string
}

type plexStats struct {
	version     string
	users       map[string]int
	modes       map[string]int
	transcoding int
	speed       float64
}

func NewCollector(url, username, password string, _ time.Duration) prometheus.Collector {
	return &Collector{
		PlexAPI: &mediaclient.PlexClient{
			Client:   &http.Client{},
			URL:      url,
			UserName: username,
			Password: password,
			Options: mediaclient.PlexOpts{
				PrometheusSummary: metrics.RequestDuration,
			},
		},
		url: url,
	}
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- transcodingMetric
	ch <- speedMetric
	ch <- usersMetric
	ch <- modesMetric
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := coll.getStats()
	if err != nil {
		log.WithError(err).Warning("failed to collect plex metrics")
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), stats.version, coll.url)
	ch <- prometheus.MustNewConstMetric(transcodingMetric, prometheus.GaugeValue, float64(stats.transcoding), coll.url)
	ch <- prometheus.MustNewConstMetric(speedMetric, prometheus.GaugeValue, stats.speed, coll.url)
	for user, count := range stats.users {
		ch <- prometheus.MustNewConstMetric(usersMetric, prometheus.GaugeValue, float64(count), user, coll.url)
	}
	for mode, count := range stats.modes {
		ch <- prometheus.MustNewConstMetric(modesMetric, prometheus.GaugeValue, float64(count), mode, coll.url)
	}
}

func (coll *Collector) getStats() (stats plexStats, err error) {
	ctx := context.Background()
	stats.version, err = coll.GetVersion(ctx)
	if err == nil {
		stats.users, stats.modes, stats.transcoding, stats.speed, err = coll.GetSessions(ctx)
	}
	return stats, err
}
