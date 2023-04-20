package plex

import (
	"context"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/v2/pkg/iplocator"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
	"strconv"
	"time"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	API
	IPLocator
	url       string
	transport *httpclient.RoundTripper
}

//go:generate mockery --name API
type API interface {
	GetIdentity(context.Context) (plex.Identity, error)
	GetSessions(context.Context) (plex.Sessions, error)
}

//go:generate mockery --name IPLocator
type IPLocator interface {
	Locate(string) (float64, float64, error)
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
	r := httpclient.NewRoundTripper(httpclient.WithMetrics("mediamon", "", "plex"))
	l := iplocator.New()
	l.Logger = slog.Default()
	return &Collector{
		API: &plex.Client{
			HTTPClient: &http.Client{Transport: r},
			URL:        url,
			UserName:   username,
			Password:   password,
		},
		IPLocator: l,
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
	start := time.Now()
	coll.collectVersion(ch)
	coll.collectSessionStats(ch)
	coll.transport.Collect(ch)
	slog.Debug("plex stats collected", "duration", time.Since(start))
}

func (coll *Collector) collectVersion(ch chan<- prometheus.Metric) {
	identity, err := coll.API.GetIdentity(context.Background())
	if err != nil {
		//ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error","Error getting Plex version", nil, nil),err)
		slog.Error("failed to collect Plex version", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.Version, coll.url)
}

func (coll *Collector) collectSessionStats(ch chan<- prometheus.Metric) {
	sessions, err := coll.API.GetSessions(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error", "Error getting Plex session stats", nil, nil), err)
		slog.Error("failed to collect Plex session stats", "err", err)
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
	if lon, lat, err := coll.IPLocator.Locate(address); err == nil {
		lonAsString = strconv.FormatFloat(lon, 'f', 2, 64)
		latAsString = strconv.FormatFloat(lat, 'f', 2, 64)
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

func parseSessions(input plex.Sessions) map[string]plexSession {
	output := make(map[string]plexSession)

	for _, session := range input.Metadata {
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
	return output
}
