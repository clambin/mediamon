package plex

import "github.com/prometheus/client_golang/prometheus"

var (
	versionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "version"),
		"version info",
		[]string{"version", "url"},
		nil,
	)

	sessionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "session_count"),
		"Active Plex session",
		[]string{"url", "id", "user", "player", "title", "location", "address", "lon", "lat"},
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
