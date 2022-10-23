package xxxarr

import "github.com/prometheus/client_golang/prometheus"

func createMetrics(application, url string) map[string]*prometheus.Desc {
	constLabels := prometheus.Labels{
		"application": application,
		"url":         url,
	}
	return map[string]*prometheus.Desc{
		"version": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "version"),
			"Version info",
			[]string{"version"},
			constLabels,
		),
		"health": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "health"),
			"Server health",
			[]string{"type"},
			constLabels,
		),
		"calendar": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "calendar"),
			"Upcoming episodes / movies",
			[]string{"title"},
			constLabels,
		),
		"queued_count": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
			"Episodes / movies being downloaded",
			nil,
			constLabels,
		),
		"queued_total": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_total_bytes"),
			"Size of episode / movie being downloaded in bytes",
			[]string{"title"},
			constLabels,
		),
		"queued_downloaded": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_downloaded_bytes"),
			"Downloaded size of episode / movie being downloaded in bytes",
			[]string{"title"},
			constLabels,
		),
		"monitored": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
			"Number of Monitored series / movies",
			nil,
			constLabels,
		),
		"unmonitored": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
			"Number of Unmonitored series / movies",
			nil,
			constLabels,
		),
	}
}
