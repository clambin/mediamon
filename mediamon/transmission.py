import logging
import requests
from pimetrics.probe import APIProbe
from mediamon import metrics


class TransmissionProbe(APIProbe):
    def __init__(self, host):
        super().__init__(f'http://{host}/')
        self.api_key = ''
        self.healthy = True

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

    def _call(self, method):
        try:
            headers = {'X-Transmission-Session-Id': self.api_key}
            body = {"method": method}
            response = self.post(endpoint='transmission/rpc', headers=headers, body=body)
            if response.status_code == 200:
                return response.json()
            elif response.status_code == 409:
                try:
                    self.api_key = response.headers['X-Transmission-Session-Id']
                    return self._call(method)
                except KeyError:
                    logging.warning('Could not get new X-Transmission-Session-Id')
            else:
                logging.warning(f'Transmission call failed: {response.status_code} - {response.reason}')
        except requests.exceptions.RequestException as err:
            logging.warning(f'Transmission call failed: {err}')
        return None

    def call(self, method):
        if response := self._call(method):
            if 'arguments' in response:
                if not self.healthy:
                    logging.info('Connection with Transmission re-established')
                    self.healthy = True
                return response['arguments']
            logging.warning('Could not parse Transmission response: missing \'arguments\' payload')
        self.healthy = False
        return None

    def measure_stats(self):
        stats = self.call('session-stats')
        return {
            'active_torrent_count': stats['activeTorrentCount'] if stats else 0,
            'paused_torrent_count': stats['pausedTorrentCount'] if stats else 0,
            'download_speed': stats['downloadSpeed'] if stats else 0,
            'upload_speed': stats['uploadSpeed'] if stats else 0,
        }

    def measure_version(self):
        stats = self.call('session-get')
        return stats['version'] if stats else 0

    def measure(self):
        return {
            'stats': self.measure_stats(),
            'version': self.measure_version(),
        }
