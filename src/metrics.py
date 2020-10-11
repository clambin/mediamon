from prometheus_client import Gauge

GAUGES = {
    'plex_session_count':
        Gauge('mediaserver_plex_session_count', 'Active Plex sessions', ['server', 'user']),
    'plex_transcoder_count':
        Gauge('mediaserver_plex_transcoder_count', 'Active Transcoder count', ['server']),
    'plex_transcoder_type_count':
        Gauge('mediaserver_plex_transcoder_type_count', 'Active Transcoder count by type', ['server', 'mode']),
    'plex_transcoder_speed_total':
        Gauge('mediaserver_plex_transcoder_speed_total', 'Speed of active transcoders', ['server']),
    'plex_transcoder_encoding_count':
        Gauge('mediaserver_plex_transcoder_encoding_count', 'Number of transcoders that are acticely encoding',
              ['server']),
    'xxxarr_calendar': Gauge('mediaserver_calendar_count', 'Number of upcoming episodes', ['server']),
    'xxxarr_queue': Gauge('mediaserver_queued_count', 'Number of queued torrents', ['server']),
    'xxxarr_monitored': Gauge('mediaserver_monitored_count', 'Number of monitored entries', ['server']),
    'xxxarr_unmonitored': Gauge('mediaserver_unmonitored_count', 'Number of unmonitored entries', ['server']),
    'version': Gauge('mediaserver_server_info', 'Server info', ['server', 'version'])
}


def report(metrics, application):
    for key, value in metrics.items():
        if key == 'plex_transcoder_type_count':
            for mode in value.keys():
                GAUGES[key].labels(application, mode).set(value[mode])
        elif key == 'plex_session_count':
            for user in value.keys():
                GAUGES[key].labels(application, user).set(value[user])
        elif key == 'version':
            GAUGES[key].labels(application, value).set(1)
        else:
            GAUGES[key].labels(application).set(value)
