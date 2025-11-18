package plex

import (
	"context"
	"iter"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"codeberg.org/clambin/go-common/set"
	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	sessionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_count"),
		"Active Plex session progress",
		[]string{"url", "user", "player", "title", "mode", "location", "address", "lon", "lat", "videoCodec", "audioCodec"},
		nil,
	)

	bandwidthMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_bandwidth"),
		"Active Plex session Bandwidth usage (in kbps)",
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

type sessionCollector struct {
	sessionGetter sessionGetter
	ipLocator     IPLocator
	logger        *slog.Logger
	url           string
}

type sessionGetter interface {
	GetSessions(context.Context) ([]plex.Session, error)
}

func (c sessionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sessionMetric
	ch <- bandwidthMetric
	ch <- transcodersMetric
	ch <- speedMetric
}

func (c sessionCollector) Collect(ch chan<- prometheus.Metric) {
	sessions, err := c.sessionGetter.GetSessions(context.Background())
	if err != nil {
		c.logger.Error("fail to collect session metrics", "err", err)
	}

	var active, throttled, speed float64

	for _, stats := range c.plexSessions(sessions) {
		ch <- prometheus.MustNewConstMetric(sessionMetric, prometheus.GaugeValue, stats.progress,
			c.url, stats.user, stats.player, stats.title, stats.videoMode, stats.location, stats.address, stats.longitude, stats.latitude, stats.videoCodec, stats.audioCodec,
		)
		ch <- prometheus.MustNewConstMetric(bandwidthMetric, prometheus.GaugeValue, float64(stats.bandwidth),
			c.url, stats.user, stats.player, stats.title, stats.videoMode, stats.location, stats.address, stats.longitude, stats.latitude, stats.videoCodec, stats.audioCodec,
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
		ch <- prometheus.MustNewConstMetric(transcodersMetric, prometheus.GaugeValue, active, c.url, "transcoding")
		ch <- prometheus.MustNewConstMetric(transcodersMetric, prometheus.GaugeValue, throttled, c.url, "throttled")
		ch <- prometheus.MustNewConstMetric(speedMetric, prometheus.GaugeValue, speed, c.url)
	}
}

func (c sessionCollector) locateAddress(address string) (lonAsString, latAsString string) {
	if location, err := c.ipLocator.Locate(address); err == nil {
		lonAsString = strconv.FormatFloat(location.Lon, 'f', 2, 64)
		latAsString = strconv.FormatFloat(location.Lat, 'f', 2, 64)
	}
	return
}

type plexSession struct {
	user       string
	player     string
	location   string
	longitude  string
	latitude   string
	title      string
	address    string
	videoMode  string
	audioCodec string
	videoCodec string
	progress   float64
	bandwidth  int
	speed      float64
	throttled  bool
}

func (c sessionCollector) plexSessions(sessions []plex.Session) iter.Seq2[string, plexSession] {
	return func(yield func(string, plexSession) bool) {
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

			s := plexSession{
				user:       session.User.Title,
				player:     session.Player.Product,
				location:   session.Session.Location,
				title:      session.GetTitle(),
				address:    session.Player.Address,
				progress:   progress,
				bandwidth:  session.Session.Bandwidth,
				videoMode:  session.GetVideoMode(),
				throttled:  session.TranscodeSession.Throttled,
				speed:      session.TranscodeSession.Speed,
				videoCodec: strings.Join(videoCodecs.List(), ","),
				audioCodec: strings.Join(audioCodecs.List(), ","),
			}

			if s.location != "lan" {
				s.longitude, s.latitude = c.locateAddress(session.Player.Address)
			}

			if !yield(session.Session.ID, s) {
				return
			}
		}
	}
}
