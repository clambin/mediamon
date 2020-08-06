import json
import logging
import requests
from prometheus_client import Gauge

from pimetrics.probe import APIProbe

GAUGES = {
    'active_torrent_count': Gauge('mediaserver_active_torrent_count', 'Active torrents'),
    'paused_torrent_count': Gauge('mediaserver_paused_torrent_count', 'Paused torrents'),
    'download_speed': Gauge('mediaserver_download_speed', 'Transmission download speed in bytes/sec'),
    'upload_speed': Gauge('mediaserver_upload_speed', 'Transmission upload speed in bytes/sec'),
}


class TransmissionProbe(APIProbe):
    def __init__(self, host):
        super().__init__(f'http://{host}/')
        self.api_key = ''
        self.connecting = True

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
                if not self.connecting:
                    logging.info('Connection with Transmission re-established')
                    self.connecting = True
                return response.json()['arguments']
            if response.status_code == 409:
                try:
                    self.api_key = response.headers['X-Transmission-Session-Id']
                    return self.measure()
                except KeyError:
                    logging.warning('Could not get new X-Transmission-Session-Id')
            else:
                logging.warning(f'Transmission call failed: {response.status_code}')
        except requests.exceptions.RequestException as err:
            logging.warning(f'Transmission call failed: {err}')
        self.connecting = False
        return None
