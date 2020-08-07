import requests
import logging
import xmltodict
import xml
from collections import OrderedDict
from pimetrics.probe import APIProbe
import src.version
from prometheus_client import Gauge

GAUGES = {
    'session_count':
        Gauge('mediaserver_plex_session_count', 'Active Plex sessions', ['server']),
    'transcoder_count':
        Gauge('mediaserver_plex_transcoder_count', 'Active Transcoder count', ['server']),
    'transcoder_type_count':
        Gauge('mediaserver_plex_transcoder_type_count', 'Active Transcoder count by type', ['server', 'mode']),
    'transcoder_speed_total':
        Gauge('mediaserver_plex_transcoder_speed_total', 'Speed of active transcoders', ['server']),
    'transcoder_encoding_count':
        Gauge('mediaserver_plex_transcoder_encoding_count', 'Number of transcoders that are acticely encoding',
              ['server']),
}


class AddressManager:
    def __init__(self, addresses):
        self.addresses = addresses
        self.address_index = 0
        self.status = None

    @property
    def address(self):
        return self.addresses[self.address_index]

    def switch(self):
        self.address_index = (self.address_index + 1) % len(self.addresses)

    @property
    def connecting(self):
        return self.status is None or self.status

    @connecting.setter
    def connecting(self, status):
        self.status = status


class PlexProbe(APIProbe, AddressManager):
    def __init__(self, authtoken, name, addresses):
        AddressManager.__init__(self, addresses)
        APIProbe.__init__(self, '')
        self.name = name
        self.headers = {
            'X-Plex-Token': authtoken,
            'Accept': 'application/json'
        }
        self.modes = set()

    def call(self, endpoint):
        first_server = None
        while self.address != first_server:
            if first_server is None:
                first_server = self.address
            url = f'{self.address}{endpoint}'
            try:
                response = requests.get(url, headers=self.headers)
                if response.status_code == 200:
                    self.connecting = True
                    return response.json()
                logging.warning(f'{self.name}: received {response.status_code} - {response.reason}')
            except requests.exceptions.ConnectionError as e:
                logging.warning(f'{self.name}: failed to connect. {e}')
            logging.warning(f'{self.name}: failed on {self.address}. Moving to next server')
            self.switch()
        logging.warning(f'{self.name}: no working servers found')
        self.connecting = False
        return None

    def report(self, output):
        logging.debug(f'Reporting {output}')
        for key, value in output.items():
            if key == 'transcoder_type_count':
                # run through all discovered modes so we report zero when a mode is no longer running
                for mode in self.modes:
                    GAUGES[key].labels(self.name, mode).set(value[mode] if mode in value else 0)
            else:
                GAUGES[key].labels(self.name).set(value)

    def process(self, output):
        logging.debug(f'Processing {output}')
        modes = set([entry['mode'] for entry in output if entry['transcode']])
        self.modes.update(modes)
        return {
            'session_count':
                len(output),
            'transcoder_count':
                len([entry for entry in output if entry['transcode']]),
            'transcoder_type_count': {
                mode: len([entry for entry in output if entry['mode'] == mode])
                for mode in modes},
            'transcoder_speed_total':
                sum([entry['speed'] for entry in output]),
            'transcoder_encoding_count':
                len([entry for entry in output if entry['transcode'] and not entry['throttled']]),
        }

    @staticmethod
    def parse_session(session):
        if 'TranscodeSession' in session:
            return {
                'transcode': True,
                'mode': session['TranscodeSession']['videoDecision'],
                'throttled': session['TranscodeSession']['throttled'],
                'speed': float(session['TranscodeSession']['speed'])
            }
        else:
            return {
                'transcode': False,
                'mode': None,
                'throttled': False,
                'speed': 0
            }

    @staticmethod
    def parse_sessions(response):
        logging.debug(f'/status/session result: {response}')
        try:
            if 'Metadata' in response['MediaContainer']:
                return [PlexProbe.parse_session(session) for session in response['MediaContainer']['Metadata']]
        except KeyError as e:
            logging.warning(f'Failed to get sessions: missing {e}')
        except TypeError as e:
            logging.warning(f'Failed to get sessions: missing {e}')
        return []

    def measure_sessions(self):
        return PlexProbe.parse_sessions(self.call('/status/sessions'))

    def measure(self):
        return self.measure_sessions()


class PlexServer:
    def __init__(self, username, password):
        self.authtext = f'user%5Blogin%5D={username}&user%5Bpassword%5D={password}'
        self.base_headers = {
            'X-Plex-Product': 'mediamon',
            'X-Plex-Version': src.version.version,
            # FIXME: generate UUID
            'X-Plex-Client-Identifier': f'mediamon-v{src.version.version}'
        }
        self.authtoken = None
        self.probes = []

    def _login(self):
        try:
            response = requests.post('https://plex.tv/users/sign_in.xml',
                                     headers=self.base_headers,
                                     data=self.authtext)
            if response.status_code == 201:
                try:
                    result = xmltodict.parse(response.content, response.encoding)
                    self.authtoken = result['user']['@authenticationToken']
                    return True
                except KeyError as e:
                    logging.error(f'Could not parse login response: {e}')
            else:
                logging.error(f'Failed to log in to plex.tv: {response.status_code} - {response.reason}')
        except requests.exceptions.ConnectionError as e:
            logging.warning(f'Failed to connect to plex.tv: {e}')
        return False

    @staticmethod
    def _parse_servers(output, encoding='UTF-8'):
        try:
            result = xmltodict.parse(output, encoding)
            return [{
                'name': device['@name'],
                'addresses':
                    [device['Connection']['@uri']] if type(device['Connection']) == OrderedDict else
                    [connection['@uri'] for connection in device['Connection']]
            } for device in result['MediaContainer']['Device'] if device['@provides'] == 'server']
        except KeyError as e:
            logging.warning(f'Failed to parse server list: missing tag {e}')
        except TypeError:
            logging.warning('Failed to parse server list. Unexpected tags found')
        except xml.parsers.expat.ExpatError as e:
            logging.warning(f'Failed to parse server list: {e}')
        return []

    def call(self, url, headers):
        # separate method so we can stub the API in unittests
        try:
            response = requests.get(url, headers=headers)
            if response.status_code == 200:
                logging.debug(response.content)
                return response.content
            else:
                logging.error(f'Failed to get server list from plex.tv: {response.status_code} - {response.reason}')
        except requests.exceptions.ConnectionError as e:
            logging.warning(f'Failed to connect to {url}: {e}')
        return None

    def _get_servers(self):
        if not self.authtoken and not self._login():
            return []
        headers = self.base_headers
        headers['X-Plex-Token'] = self.authtoken
        response = self.call('https://plex.tv/devices.xml', headers=headers)
        if response:
            return self._parse_servers(response)
        return []

    def make_probes(self):
        servers = self._get_servers()
        for server in servers:
            logging.info(f'Plex server found: {server["name"]}: {",".join(server["addresses"])}')
        self.probes = [PlexProbe(self.authtoken, server['name'], server['addresses']) for server in servers]
        return self.probes

    def _healthcheck(self):
        unhealthy = [probe for probe in self.probes if probe.connecting is False]
        if unhealthy:
            healthy = [probe for probe in self.probes if probe.connecting]
            # TODO: force log in again to get a fresh authtoken?
            servers = self._get_servers()
            for probe in unhealthy:
                for server in servers:
                    if server['name'] == probe.name:
                        logging.info(f'Reconnecting to {probe.name}')
                        healthy.append(PlexProbe(self.authtoken, server['name'], server['addresses']))
            self.probes = healthy

    def run(self):
        if not self.probes:
            self.make_probes()
        for probe in self.probes:
            probe.run()
        self._healthcheck()
