package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"mediamon/pkg/metrics"
)

var (
	MediaServerVersion = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_server_info",
		Help: "Server info",
	}, []string{"server", "version"})

	PlexSessionCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_plex_session_count",
		Help: "Active Plex sessions",
	}, []string{"user"})

	PlexTranscoderTypeCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_plex_transcoder_type_count",
		Help: "Active Transcoder count by type",
	}, []string{"mode"})

	XXXArrCalendarCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_calendar_count",
		Help: "Number of upcoming episodes / movies",
	}, []string{"server"})

	XXXarrQueuedCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_queued_count",
		Help: "Number of queued torrents",
	}, []string{"server"})

	XXXarrMonitoredCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_monitored_count",
		Help: "Number of monitored series / movies",
	}, []string{"server"})

	XXXarrUnmonitoredCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_unmonitored_count",
		Help: "Number of unmonitored series / movies",
	}, []string{"server"})

	PlexTranscoderSpeedTotal = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_plex_transcoder_speed_total",
		Help: "Speed of active transcoders",
	})
	PlexTranscoderEncodingCount = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_plex_transcoder_encoding_count",
		Help: "Number of active transcoders",
	})
	TransmissionActiveTorrentCount = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_active_torrent_count",
		Help: "Number of active torrents",
	})

	TransmissionPausedTorrentCount = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_paused_torrent_count",
		Help: "Number of paused torrents",
	})

	TransmissionDownloadSpeed = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_download_speed",
		Help: "Transmission download speed in bytes / sec",
	})

	TransmissionUploadSpeed = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_upload_speed",
		Help: "Transmission upload speed in bytes / sec",
	})

	OpenVPNClientStatus = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "openvpn_client_status",
		Help: "OpenVPN Client Status",
	})

	OpenVPNClientReadTotal = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "openvpn_client_tcp_udp_read_bytes_total",
		Help: "OpenVPN client bytes read",
	})

	OpenVPNClientWriteTotal = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "openvpn_client_tcp_udp_write_bytes_total",
		Help: "OpenVPN client bytes written",
	})
)
