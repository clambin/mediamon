package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	labeledGauges = map[string]*prometheus.GaugeVec{
		"version": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_server_info",
			Help: "Server info",
		}, []string{"server", "version"}),
		"plex_session_count": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_plex_session_count",
			Help: "Active Plex sessions",
		}, []string{"user"}),
		"plex_transcoder_type_count": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_plex_transcoder_type_count",
			Help: "Active Transcoder count by type",
		}, []string{"mode"}),
		"xxxarr_calendar": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_calendar_count",
			Help: "Number of upcoming episodes / movies",
		}, []string{"server"}),
		"xxxarr_queue": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_queue_count",
			Help: "Number of queued torrents",
		}, []string{"server"}),
		"xxxarr_monitored": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_monitored_count",
			Help: "Number of monitored series / movies",
		}, []string{"server"}),
		"xxxarr_unmonitored": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_unmonitored_count",
			Help: "Number of unmonitored series / movies",
		}, []string{"server"}),
	}

	unlabeledGauges = map[string]prometheus.Gauge{
		"plex_transcoder_speed_total": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "mediaserver_plex_transcoder_speed_total",
			Help: "Speed of active transcoders",
		}),
		"plex_transcoder_encoding_count": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "mediaserver_plex_transcoder_encoding_count",
			Help: "Number of active transcoders",
		}),
		"active_torrent_count": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "mediaserver_active_torrent_count",
			Help: "Number of active torrents",
		}),
		"paused_torrent_count": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "mediaserver_paused_torrent_count",
			Help: "Number of paused torrents",
		}),
		"download_speed": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "mediaserver_download_speed",
			Help: "Transmission download speed in bytes / sec",
		}),
		"upload_speed": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "mediaserver_upload_speed",
			Help: "Transmission upload speed in bytes / sec",
		}),
	}

	cachedValues = map[string]map[string]float64{}
)

// Init initializes the prometheus metrics server
func Init(port int) {
	http.Handle("/metrics", promhttp.Handler())
	listenAddress := fmt.Sprintf(":%d", port)
	go func(listenAddr string) {
		err := http.ListenAndServe(listenAddress, nil)
		if err != nil {
			panic(err)
		}
	}(listenAddress)
}

// Publish pushes the specified metric to Prometheus
func Publish(metric string, value float64, labels ...string) bool {
	log.Debugf("%s(%s): %f", metric, labels, value)
	if gauge, ok := unlabeledGauges[metric]; ok {
		gauge.Set(value)
		SaveValue(metric, value, labels...)
		return true
	} else if gauge, ok := labeledGauges[metric]; ok {
		gauge.WithLabelValues(labels...).Set(value)
		SaveValue(metric, value, labels...)
		return true
	}
	log.Warningf("metric '%s' not found", metric)
	return false
}

// SaveValue stores the last value reported so unit tests can verify the correct value was reported
func SaveValue(metric string, value float64, labels ...string) {
	subMap, ok := cachedValues[metric]
	if ok == false {
		subMap = make(map[string]float64)
		cachedValues[metric] = subMap
	}
	key := strings.Join(labels, "|")
	subMap[key] = value
}

// LoadValue gets the last value reported so unit tests can verify the correct value was reported
func LoadValue(metric string, labels ...string) (float64, bool) {
	if value, ok := cachedValues[metric][strings.Join(labels, "|")]; ok {
		return value, true
	}
	return 0, false
}
