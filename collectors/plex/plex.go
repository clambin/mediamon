package plex

import (
	"context"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
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

	transcodersMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_total_count"),
		"Number of transcode sessions",
		[]string{"url"},
		nil,
	)

	transcodingMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_active_count"),
		"Number of active transcode sessions",
		[]string{"url"},
		nil,
	)

	speedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_speed_total"),
		"Total speed of active transcoders",
		[]string{"url"},
		nil,
	)

	locationMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_location_count"),
		"Active plex sessions by user",
		[]string{"location", "url"},
		nil,
	)

	usersMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_user_count"),
		"Active Plex sessions by location",
		[]string{"user", "url"},
		nil,
	)
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	plex.API
	url string
}

// NewCollector creates a new Collector
func NewCollector(url, username, password string, _ time.Duration) prometheus.Collector {
	return &Collector{
		API: &plex.Client{
			Client:   &http.Client{},
			URL:      url,
			UserName: username,
			Password: password,
			Options: plex.Options{
				PrometheusSummary: metrics.RequestDuration,
			},
		},
		url: url,
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- transcodersMetric
	ch <- transcodingMetric
	ch <- speedMetric
	ch <- locationMetric
	ch <- usersMetric
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	coll.collectVersion(ch)
	coll.collectSessionStats(ch)
}

func (coll *Collector) collectVersion(ch chan<- prometheus.Metric) {
	version, err := coll.GetVersion(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("mediamon_error",
				"Error getting Plex version", nil, nil),
			err)
		log.WithError(err).Warning("failed to collect Plex version")
		return
	}

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), version, coll.url)
}

func (coll *Collector) collectSessionStats(ch chan<- prometheus.Metric) {
	sessions, err := coll.GetSessions(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("mediamon_error",
				"Error getting Plex session stats", nil, nil),
			err)
		log.WithError(err).Warning("failed to collect Plex session stats")
		return
	}

	users := map[string]int{}
	transcoders := 0.0
	transcoding := 0.0
	speed := 0.0
	lan := 0.0
	wan := 0.0

	for _, session := range sessions {
		if session.Transcode {
			transcoders++
			if session.Throttled == false {
				transcoding++
			}
			speed += session.Speed
		}
		if session.Local {
			lan++
		} else {
			wan++
		}
		count, _ := users[session.User]
		count++
		users[session.User] = count
	}

	ch <- prometheus.MustNewConstMetric(transcodersMetric, prometheus.GaugeValue, transcoders, coll.url)
	ch <- prometheus.MustNewConstMetric(transcodingMetric, prometheus.GaugeValue, transcoding, coll.url)
	ch <- prometheus.MustNewConstMetric(speedMetric, prometheus.GaugeValue, speed, coll.url)
	ch <- prometheus.MustNewConstMetric(locationMetric, prometheus.GaugeValue, lan, "lan", coll.url)
	ch <- prometheus.MustNewConstMetric(locationMetric, prometheus.GaugeValue, wan, "wan", coll.url)

	for user, count := range users {
		ch <- prometheus.MustNewConstMetric(usersMetric, prometheus.GaugeValue, float64(count), user, coll.url)
	}
}
