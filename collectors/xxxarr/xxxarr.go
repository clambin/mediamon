package xxxarr

import (
	metrics2 "github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/cache"
	"github.com/clambin/mediamon/collectors/xxxarr/updater"
	"github.com/clambin/mediamon/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// Collector presents Sonarr/Radarr statistics as Prometheus metrics
type Collector struct {
	updater.StatsGetter
	cache.Cache[updater.Stats]
	application string
	metrics     map[string]*prometheus.Desc
}

// NewRadarrCollector creates a new RadarrCollector
func NewRadarrCollector(url, apiKey string, interval time.Duration) prometheus.Collector {
	getter := updater.RadarrUpdater{
		Client: xxxarr.NewRadarrClient(apiKey, url, xxxarr.Options{
			PrometheusMetrics: metrics2.APIClientMetrics{
				Latency: metrics.Latency,
				Errors:  metrics.Errors,
			},
		}),
	}

	return &Collector{
		StatsGetter: getter,
		Cache: cache.Cache[updater.Stats]{
			Duration: interval,
			Updater:  getter.GetStats,
		},
		application: "radarr",
		metrics:     createMetrics("radarr"),
	}
}

// NewSonarrCollector creates a new SonarrCollector
func NewSonarrCollector(url, apiKey string, interval time.Duration) prometheus.Collector {
	getter := &updater.SonarrUpdater{
		Client: xxxarr.NewSonarrClient(apiKey, url, xxxarr.Options{
			PrometheusMetrics: metrics2.APIClientMetrics{
				Latency: metrics.Latency,
				Errors:  metrics.Errors,
			},
		}),
	}

	return &Collector{
		StatsGetter: getter,
		Cache: cache.Cache[updater.Stats]{
			Duration: interval,
			Updater:  getter.GetStats,
		},
		application: "sonarr",
		metrics:     createMetrics("sonarr"),
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range coll.metrics {
		ch <- metric
	}
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Cache.Update()

	ch <- prometheus.MustNewConstMetric(coll.metrics["version"], prometheus.GaugeValue, float64(1), stats.Version, stats.URL)
	for _, title := range stats.Calendar {
		ch <- prometheus.MustNewConstMetric(coll.metrics["calendar"], prometheus.GaugeValue, 1.0, stats.URL, title)
	}

	ch <- prometheus.MustNewConstMetric(coll.metrics["queued_count"], prometheus.GaugeValue, float64(len(stats.Queued)), stats.URL)
	for _, queued := range stats.Queued {
		ch <- prometheus.MustNewConstMetric(coll.metrics["queued_total"], prometheus.GaugeValue, queued.TotalBytes, stats.URL, queued.Name)
		ch <- prometheus.MustNewConstMetric(coll.metrics["queued_downloaded"], prometheus.GaugeValue, queued.DownloadedBytes, stats.URL, queued.Name)
	}
	ch <- prometheus.MustNewConstMetric(coll.metrics["monitored"], prometheus.GaugeValue, float64(stats.Monitored), stats.URL)
	ch <- prometheus.MustNewConstMetric(coll.metrics["unmonitored"], prometheus.GaugeValue, float64(stats.Unmonitored), stats.URL)
}
