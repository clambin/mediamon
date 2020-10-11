import json
import logging
import requests
from prometheus_client import Gauge
from pimetrics.probe import APIProbe
from src import metrics


class TransmissionProbe(APIProbe):
    def __init__(self, host):
        super().__init__(f'http://{host}/')
        self.api_key = ''
        self.connecting = True

    def report(self, output):
        metrics.report(output, 'transmission')

    def process(self, output):
        return {
            'active_torrent_count': output['stats']['active_torrent_count'],
            'paused_torrent_count': output['stats']['paused_torrent_count'],
            'download_speed': output['stats']['download_speed'],
            'upload_speed': output['stats']['upload_speed'],
            'version': output['version'],
        }

    def call(self, method):
        try:
            headers = {'X-Transmission-Session-Id': self.api_key}
            body = {"method": method}
            response = self.post(endpoint='transmission/rpc', headers=headers, body=body)
            if response.status_code == 200:
                if not self.connecting:
                    logging.info('Connection with Transmission re-established')
                    self.connecting = True
                return response.json()['arguments']
            if response.status_code == 409:
                try:
                    self.api_key = response.headers['X-Transmission-Session-Id']
                    return self.call(method)
                except KeyError:
                    logging.warning('Could not get new X-Transmission-Session-Id')
            else:
                logging.warning(f'Transmission call failed: {response.status_code}')
        except requests.exceptions.RequestException as err:
            logging.warning(f'Transmission call failed: {err}')
        self.connecting = False
        return None

    def measure_stats(self):
        stats = self.call('session-stats')
        return {
            'active_torrent_count': stats['activeTorrentCount'],
            'paused_torrent_count': stats['pausedTorrentCount'],
            'download_speed': stats['downloadSpeed'],
            'upload_speed': stats['uploadSpeed'],
        }

    def measure_version(self):
        stats = self.call('session-get')
        return stats['version']

    def measure(self):
        return {
            'stats': self.measure_stats(),
            'version': self.measure_version(),
        }
