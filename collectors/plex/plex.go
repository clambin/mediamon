package plex

import (
	"context"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediamon/v2/pkg/iplocator"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"strconv"
	"time"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	API
	IPLocator
	url       string
	transport *httpclient.RoundTripper
	logger    *slog.Logger
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
func NewCollector(version, url, username, password string) *Collector {
	r := httpclient.NewRoundTripper(httpclient.WithMetrics("mediamon", "", "plex"))
	return &Collector{
		API:       plex.New(username, password, "github.com/clambin/mediamon", version, url, r),
		IPLocator: iplocator.New(),
		url:       url,
		transport: r,
		logger:    slog.Default().With("collector", "plex"),
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
	coll.logger.Debug("stats collected", "duration", time.Since(start))
}

func (coll *Collector) collectVersion(ch chan<- prometheus.Metric) {
	identity, err := coll.API.GetIdentity(context.Background())
	if err != nil {
		//ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error","Error getting Plex version", nil, nil),err)
		coll.logger.Error("failed to collect version", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.Version, coll.url)
}

func (coll *Collector) collectSessionStats(ch chan<- prometheus.Metric) {
	sessions, err := coll.API.GetSessions(context.Background())
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error", "Error getting Plex session stats", nil, nil), err)
		coll.logger.Error("failed to collect session stats", "err", err)
		return
	}

	for _, session := range sessions.Metadata {
		coll.log(session)
	}

	var active, throttled, speed float64

	for id, stats := range parseSessions(sessions) {
		var lon, lat string
		if stats.location != "lan" {
			lon, lat = coll.locateAddress(stats.address)
		}

		ch <- prometheus.MustNewConstMetric(sessionMetric, prometheus.GaugeValue, stats.progress,
			coll.url, id, stats.user, stats.player, stats.title, stats.videoMode, stats.location, stats.address, lon, lat,
		)

		if stats.videoMode == "transcode" {
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
	progress  float64
	videoMode string
	throttled bool
	speed     float64
}

func parseSessions(input plex.Sessions) map[string]plexSession {
	output := make(map[string]plexSession)

	for _, session := range input.Metadata {
		output[session.Session.ID] = plexSession{
			user:      session.User.Title,
			player:    session.Player.Product,
			location:  session.Session.Location,
			title:     session.GetTitle(),
			address:   session.Player.Address,
			progress:  session.GetProgress(),
			videoMode: session.GetVideoMode(),
			throttled: session.TranscodeSession.Throttled,
			speed:     session.TranscodeSession.Speed,
		}
	}
	return output
}

type sessionInfo struct {
	parts []partInfo
}

type partInfo struct {
	decision   string
	streamInfo []streamInfo
}

type streamInfo struct {
	decision string
	kind     int
	codec    string
	location string
}

func (coll *Collector) log(session plex.Session) {
	var sessionInfos []sessionInfo
	for _, media := range session.Media {
		var parts []partInfo
		for _, part := range media.Part {
			var streams []streamInfo
			for _, stream := range part.Stream {
				streams = append(streams, streamInfo{
					decision: stream.Decision,
					kind:     stream.StreamType,
					codec:    stream.Codec,
					location: stream.Location,
				})
			}
			parts = append(parts, partInfo{
				decision:   part.Decision,
				streamInfo: streams,
			})
		}
		sessionInfos = append(sessionInfos, sessionInfo{
			parts: parts,
		})
	}

	coll.logger.Debug("session found",
		"title", session.GetTitle(),
		"mode", session.GetVideoMode(),
		"media.part.decisions", sessionInfos,
		"transcode", session.TranscodeSession,
	)
}
