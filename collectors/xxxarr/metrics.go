package xxxarr

import "github.com/prometheus/client_golang/prometheus"

func createMetrics(application string) map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"version": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "version"),
			"Version info",
			[]string{"version", "url"},
			prometheus.Labels{"application": application},
		),
		"calendar": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "calendar"),
			"Upcoming episodes / movies",
			[]string{"url", "title"},
			prometheus.Labels{"application": application},
		),
		"queued_count": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_count"),
			"Episodes / movies being downloaded",
			[]string{"url"},
			prometheus.Labels{"application": application},
		),
		"queued_total": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_total_bytes"),
			"Size of episode / movie being downloaded in bytes",
			[]string{"url", "title"},
			prometheus.Labels{"application": application},
		),
		"queued_downloaded": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "queued_downloaded_bytes"),
			"Downloaded size of episode / movie being downloaded in bytes",
			[]string{"url", "title"},
			prometheus.Labels{"application": application},
		),
		"monitored": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "monitored_count"),
			"Number of Monitored series / movies",
			[]string{"url"},
			prometheus.Labels{"application": application},
		),
		"unmonitored": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "xxxarr", "unmonitored_count"),
			"Number of Unmonitored series / movies",
			[]string{"url"},
			prometheus.Labels{"application": application},
		),
	}
}
