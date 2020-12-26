package metrics

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
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
		"xxxarr_queued": promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mediaserver_queued_count",
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
		"openvpn_client_status": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "openvpn_client_status",
			Help: "OpenVPN Client Status",
		}),
		"openvpn_client_tcp_udp_read_bytes_total": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "openvpn_client_tcp_udp_read_bytes_total",
			Help: "OpenVPN client bytes read",
		}),
		"openvpn_client_tcp_udp_write_bytes_total": promauto.NewGauge(prometheus.GaugeOpts{
			Name: "openvpn_client_tcp_udp_write_bytes_total",
			Help: "OpenVPN client bytes written",
		}),
	}
)

// Run initialized & runs the metrics
func Run(port int, background bool) {
	listenAddress := fmt.Sprintf(":%d", port)
	http.Handle("/metrics", promhttp.Handler())

	listenFunc := func(listenAddr string) {
		err := http.ListenAndServe(listenAddress, nil)
		if err != nil {
			panic(err)
		}
	}

	if background {
		go listenFunc(listenAddress)
	} else {
		listenFunc(listenAddress)
	}
}

// Publish pushes the specified metric to Prometheus
func Publish(metric string, value float64, labels ...string) bool {
	log.Debugf("%s(%s): %f", metric, labels, value)
	if gauge, ok := unlabeledGauges[metric]; ok {
		gauge.Set(value)
		return true
	} else if gauge, ok := labeledGauges[metric]; ok {
		gauge.WithLabelValues(labels...).Set(value)
		return true
	}
	log.Warningf("metric '%s' not found", metric)
	return false
}

// LoadValue gets the last value reported so unit tests can verify the correct value was reported
func LoadValue(metric string, labels ...string) (float64, error) {
	log.Debugf("%s(%s)", metric, labels)
	if gauge, ok := unlabeledGauges[metric]; ok {
		var m = &dto.Metric{}
		_ = gauge.Write(m)
		return m.Gauge.GetValue(), nil
	} else if gauge, ok := labeledGauges[metric]; ok {
		var m = &dto.Metric{}
		_ = gauge.WithLabelValues(labels...).Write(m)
		return m.Gauge.GetValue(), nil
	}
	return 0, errors.New("could not find " + metric)
}
