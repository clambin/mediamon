package metrics

import (
	"github.com/clambin/gotools/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// MediaServerVersion is always 1.  The "version" label contains the current version
	MediaServerVersion = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_server_info",
		Help: "Server info",
	}, []string{"server", "version"})

	// PlexSessionCount Active Plex sessions
	PlexSessionCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_plex_session_count",
		Help: "Active Plex sessions",
	}, []string{"user"})

	// PlexTranscoderTypeCount Active Transcoder count by type
	PlexTranscoderTypeCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_plex_transcoder_type_count",
		Help: "Active Transcoder count by type",
	}, []string{"mode"})

	// XXXArrCalendarCount number of episodes / movies aired today and tomorrow
	XXXArrCalendarCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_calendar_count",
		Help: "Number of upcoming episodes / movies",
	}, []string{"server"})

	// XXXArrQueuedCount number of episodes / movies being downloaded
	XXXArrQueuedCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_queued_count",
		Help: "Number of episodes / movies being downloaded",
	}, []string{"server"})

	// XXXArrMonitoredCount number of monitored series / movies
	XXXArrMonitoredCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_monitored_count",
		Help: "Number of monitored series / movies",
	}, []string{"server"})

	// XXXArrUnmonitoredCount number of unmonitored series / movies
	XXXArrUnmonitoredCount = metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mediaserver_unmonitored_count",
		Help: "Number of unmonitored series / movies",
	}, []string{"server"})

	// PlexTranscoderSpeedTotal speed of active transcoders
	PlexTranscoderSpeedTotal = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_plex_transcoder_speed_total",
		Help: "Speed of active transcoders",
	})

	// PlexTranscoderEncodingCount number of active transcoders
	PlexTranscoderEncodingCount = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_plex_transcoder_encoding_count",
		Help: "Number of active transcoders",
	})

	// TransmissionActiveTorrentCount Number of active torrents
	TransmissionActiveTorrentCount = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_active_torrent_count",
		Help: "Number of active torrents",
	})

	// TransmissionPausedTorrentCount Number of paused torrents
	TransmissionPausedTorrentCount = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_paused_torrent_count",
		Help: "Number of paused torrents",
	})

	// TransmissionDownloadSpeed Transmission download speed in bytes / sec
	TransmissionDownloadSpeed = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_download_speed",
		Help: "Transmission download speed in bytes / sec",
	})

	// TransmissionUploadSpeed Transmission upload speed in bytes / sec
	TransmissionUploadSpeed = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "mediaserver_upload_speed",
		Help: "Transmission upload speed in bytes / sec",
	})

	// OpenVPNClientStatus OpenVPN Client Status
	OpenVPNClientStatus = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "openvpn_client_status",
		Help: "OpenVPN Client Status",
	})

	// OpenVPNClientReadTotal OpenVPN client bytes read
	OpenVPNClientReadTotal = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "openvpn_client_tcp_udp_read_bytes_total",
		Help: "OpenVPN client bytes read",
	})

	// OpenVPNClientWriteTotal OpenVPN client bytes written
	OpenVPNClientWriteTotal = metrics.NewGauge(prometheus.GaugeOpts{
		Name: "openvpn_client_tcp_udp_write_bytes_total",
		Help: "OpenVPN client bytes written",
	})
)
