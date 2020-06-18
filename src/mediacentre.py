import json
import logging
from enum import Enum
import requests
from prometheus_client import Gauge

from pimetrics.probe import APIProbe

GAUGES = {
    'active_torrent_count': Gauge('mediaserver_active_torrent_count', 'Active torrents'),
    'paused_torrent_count': Gauge('mediaserver_paused_torrent_count', 'Paused torrents'),
    'download_speed': Gauge('mediaserver_download_speed', 'Transmission download speed in bytes/sec'),
    'upload_speed': Gauge('mediaserver_upload_speed', 'Transmission upload speed in bytes/sec'),
    'calendar_count': Gauge('mediaserver_calendar_count', 'Number of upcoming episodes', ['server']),
    'queued_count': Gauge('mediaserver_queued_count', 'Number of queued torrents', ['server']),
    'monitored_count': Gauge('mediaserver_monitored_count', 'Number of monitored entries', ['server']),
    'unmonitored_count': Gauge('mediaserver_unmonitored_count', 'Number of unmonitored entries', ['server']),
}


class TransmissionProbe(APIProbe):
    def __init__(self, host):
        super().__init__(f'http://{host}/')
        self.api_key = ''

    def report(self, output):
        if output:
            try:
                GAUGES['active_torrent_count'].set(output['activeTorrentCount'])
                GAUGES['paused_torrent_count'].set(output['pausedTorrentCount'])
                GAUGES['download_speed'].set(output['downloadSpeed'])
                GAUGES['upload_speed'].set(output['uploadSpeed'])
            except KeyError as err:
                logging.warning(f'Incomplete output: {err}')
                logging.debug(json.dumps(output, indent=3))

    def measure(self):
        try:
            headers = {'X-Transmission-Session-Id': self.api_key}
            body = {"method": "session-stats"}
            response = self.post(endpoint='transmission/rpc', headers=headers, body=body)
            if response.status_code == 200:
                return response.json()['arguments']
            if response.status_code == 409:
                try:
                    self.api_key = response.headers['X-Transmission-Session-Id']
                    return self.measure()
                except KeyError:
                    logging.warning('Could not get new X-Transmission-Session-Id')
                    return None
            logging.warning(f'Transmission call failed: {response.status_code}')
        except requests.exceptions.RequestException as err:
            logging.warning(f'Transmission call failed: {err}')
        return None


class MonitorProbe(APIProbe):
    class App(Enum):
        sonarr = 0
        radarr = 1

    def __init__(self, host, app, api_key):
        super().__init__(f'http://{host}/')
        self.api_key = api_key
        self.app = app

    def report(self, output):
        if output:
            calendar = output['calendar']
            queue = output['queue']
            monitored = output['monitored'][0]
            unmonitored = output['monitored'][1]
            GAUGES['calendar_count'].labels(self.app.name).set(calendar)
            GAUGES['queued_count'].labels(self.app.name).set(queue)
            GAUGES['monitored_count'].labels(self.app.name).set(monitored)
            GAUGES['unmonitored_count'].labels(self.app.name).set(unmonitored)

    def call(self, endpoint):
        result = None
        try:
            headers = {'X-Api-Key': self.api_key}
            response = self.get(endpoint=endpoint, headers=headers)
            if response.status_code == 200:
                result = response.json()
            else:
                logging.error("%d - %s" % (response.status_code, response.reason))
        except requests.exceptions.RequestException as err:
            logging.warning(f'Failed to call "{self.url}": "{err}')
        return result

    def measure_calendar(self):
        calendar = self.call('api/calendar')
        calendar = list(filter(lambda entry: not entry['hasFile'], calendar))
        return len(calendar)

    def measure_queue(self):
        queue = self.call('api/queue')
        return len(queue)

    def measure_monitored(self):
        if self.app == self.App.sonarr:
            entries = self.call('api/series')
        elif self.app == self.App.radarr:
            entries = self.call('api/movie')
        monitored = list(filter(lambda entry: entry['monitored'], entries))
        unmonitored = list(filter(lambda entry: not entry['monitored'], entries))
        return len(monitored), len(unmonitored)

    def measure(self):
        return {
            'calendar': self.measure_calendar(),
            'queue': self.measure_queue(),
            'monitored': self.measure_monitored()
        }
