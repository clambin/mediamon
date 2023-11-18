package plex

import (
	"context"
	"github.com/clambin/go-common/set"
	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"math"
	"strconv"
	"strings"
)

var (
	sessionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_count"),
		"Active Plex session",
		[]string{"url", "user", "player", "title", "mode", "location", "address", "lon", "lat", "videoCodec", "audioCodec"},
		nil,
	)

	transcodersMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_count"),
		"Video transcode session",
		[]string{"url", "state"},
		nil,
	)

	speedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "transcoder_speed"),
		"Speed of active transcoder",
		[]string{"url"},
		nil,
	)
)

var _ prometheus.Collector = sessionCollector{}

type sessionCollector struct {
	sessionGetter
	IPLocator
	url    string
	logger *slog.Logger
}

type sessionGetter interface {
	GetSessions(context.Context) ([]plex.Session, error)
}

func (c sessionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sessionMetric
	ch <- transcodersMetric
	ch <- speedMetric
}

func (c sessionCollector) Collect(ch chan<- prometheus.Metric) {
	sessions, err := c.sessionGetter.GetSessions(context.Background())
	if err != nil {
		//ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error", "Error getting Plex session stats", nil, nil), err)
		c.logger.Error("failed to collect session stats", "err", err)
		return
	}

	var active, throttled, speed float64

	for _, stats := range parseSessions(sessions) {
		var lon, lat string
		if stats.location != "lan" {
			lon, lat = c.locateAddress(stats.address)
		}

		ch <- prometheus.MustNewConstMetric(sessionMetric, prometheus.GaugeValue, stats.progress,
			c.url, stats.user, stats.player, stats.title, stats.videoMode, stats.location, stats.address, lon, lat, stats.videoCodec, stats.audioCodec,
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
			c.url, "transcoding",
		)
		ch <- prometheus.MustNewConstMetric(transcodersMetric, prometheus.GaugeValue, throttled,
			c.url, "throttled",
		)
		ch <- prometheus.MustNewConstMetric(speedMetric, prometheus.GaugeValue, speed,
			c.url,
		)
	}
}

func (c sessionCollector) locateAddress(address string) (lonAsString, latAsString string) {
	if lon, lat, err := c.IPLocator.Locate(address); err == nil {
		lonAsString = strconv.FormatFloat(lon, 'f', 2, 64)
		latAsString = strconv.FormatFloat(lat, 'f', 2, 64)
	}
	return
}

type plexSession struct {
	user       string
	player     string
	location   string
	title      string
	address    string
	progress   float64
	videoMode  string
	throttled  bool
	speed      float64
	audioCodec string
	videoCodec string
}

func parseSessions(sessions []plex.Session) map[string]plexSession {
	output := make(map[string]plexSession)

	for _, session := range sessions {
		videoCodecs := set.New[string]()
		audioCodecs := set.New[string]()
		for _, media := range session.Media {
			videoCodecs.Add(media.VideoCodec)
			audioCodecs.Add(media.AudioCodec)
		}
		progress := session.GetProgress()
		if math.IsNaN(progress) {
			progress = 0
		}

		output[session.Session.ID] = plexSession{
			user:       session.User.Title,
			player:     session.Player.Product,
			location:   session.Session.Location,
			title:      session.GetTitle(),
			address:    session.Player.Address,
			progress:   progress,
			videoMode:  session.GetVideoMode(),
			throttled:  session.TranscodeSession.Throttled,
			speed:      session.TranscodeSession.Speed,
			videoCodec: strings.Join(videoCodecs.List(), ","),
			audioCodec: strings.Join(audioCodecs.List(), ","),
		}
	}
	return output
}
