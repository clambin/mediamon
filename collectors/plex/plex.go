package plex

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/pkg/iplocator"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	plex.API
	iplocator.Locator
	url       string
	transport *httpclient.RoundTripper
}

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	URL      string
	UserName string
	Password string
}

// NewCollector creates a new Collector
func NewCollector(url, username, password string) *Collector {
	r := httpclient.NewRoundTripper(httpclient.WithRoundTripperMetrics{Namespace: "mediamon", Application: "plex"})
	return &Collector{
		API: &plex.Client{
			HTTPClient: &http.Client{Transport: r},
			URL:        url,
			UserName:   username,
			Password:   password,
		},
		Locator:   iplocator.New(),
		url:       url,
		transport: r,
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- sessionMetric
	ch <- transcodersMetric
	ch <- speedMetric
	coll.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	coll.collectVersion(ch)
	coll.collectSessionStats(ch)
	coll.transport.Collect(ch)
}

func (coll *Collector) collectVersion(ch chan<- prometheus.Metric) {
	identity, err := coll.API.GetIdentity(context.Background())
	if err != nil {
		/*
			ch <- prometheus.NewInvalidMetric(
				prometheus.NewDesc("mediamon_error",
					"Error getting Plex version", nil, nil),
				err)
			log.WithError(err).Warning("failed to collect Plex version")
		*/
		return
	}

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.MediaContainer.Version, coll.url)
}

func (coll *Collector) collectSessionStats(ch chan<- prometheus.Metric) {
	sessions, err := coll.API.GetSessions(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("mediamon_error",
				"Error getting Plex session stats", nil, nil),
			err)
		slog.Error("failed to collect Plex session stats", err)
		return
	}

	var active, throttled, speed float64

	for id, stats := range parseSessions(sessions) {
		var lon, lat string
		if stats.location != "lan" {
			lon, lat = coll.locateAddress(stats.address)
		}

		ch <- prometheus.MustNewConstMetric(sessionMetric, prometheus.GaugeValue, 1.0,
			coll.url, id, stats.user, stats.player, stats.title, stats.location, stats.address, lon, lat,
		)

		if stats.transcode {
			if stats.throttled {
				throttled++
			} else {
				active++
			}
			speed += stats.speed
		}
	}
	if active+throttled > 0 {
		ch <- prometheus.MustNewConstMetric(transcodersMetric, prometheus.GaugeValue, active,
			coll.url, "transcoding",
		)
		ch <- prometheus.MustNewConstMetric(transcodersMetric, prometheus.GaugeValue, throttled,
			coll.url, "throttled",
		)
		ch <- prometheus.MustNewConstMetric(speedMetric, prometheus.GaugeValue, speed,
			coll.url,
		)
	}
}

func (coll *Collector) locateAddress(address string) (lonAsString, latAsString string) {
	lon, lat, err := coll.Locator.Locate(address)

	if err == nil {
		lonAsString = fmt.Sprintf("%.2f", lon)
		latAsString = fmt.Sprintf("%.2f", lat)
	}
	return
}

type plexSession struct {
	user      string
	player    string
	location  string
	title     string
	address   string
	transcode bool
	throttled bool
	speed     float64
}

func parseSessions(input plex.SessionsResponse) (output map[string]plexSession) {
	output = make(map[string]plexSession)

	for _, session := range input.MediaContainer.Metadata {
		var title string
		if session.Type == "episode" {
			title = session.GrandparentTitle + " - " + session.ParentTitle + " - " + session.Title
		} else {
			title = session.Title
		}

		output[session.Session.ID] = plexSession{
			user:      session.User.Title,
			player:    session.Player.Product,
			location:  session.Session.Location,
			title:     title,
			address:   session.Player.Address,
			transcode: session.TranscodeSession.VideoDecision == "transcode",
			throttled: session.TranscodeSession.Throttled,
			speed:     session.TranscodeSession.Speed,
		}
	}
	return
}
