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
        Gauge('mediaserver_plex_session_count', 'Active Plex sessions', ['server', 'user']),
    'transcoder_count':
        Gauge('mediaserver_plex_transcoder_count', 'Active Transcoder count', ['server']),
    'transcoder_type_count':
        Gauge('mediaserver_plex_transcoder_type_count', 'Active Transcoder count by type', ['server', 'mode']),
    'transcoder_speed_total':
        Gauge('mediaserver_plex_transcoder_speed_total', 'Speed of active transcoders', ['server']),
    'transcoder_encoding_count':
        Gauge('mediaserver_plex_transcoder_encoding_count', 'Number of transcoders that are acticely encoding',
              ['server']),
    'server_info': Gauge('mediaserver_plex_info', 'Plex version', ['server', 'version']),
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
        self.users = set()
        self.modes = set()

    def call(self, endpoint):
        first_server = None
        while self.address != first_server:
            try:
                if first_server is None:
                    first_server = self.address
                url = f'{self.address}{endpoint}'
                response = requests.get(url, headers=self.headers)
                if response.status_code == 200:
                    logging.debug(response.headers)
                    if 'X-Plex-Protocol' in response.headers:
                        if self.connecting is False:
                            logging.info(f'{self.name}: connection established on {self.address}')
                            self.connecting = True
                        return response.json()
                    else:
                        logging.info(f'{url} responded, but X-Plex-Protocol header missing. Ignoring server.')
                else:
                    logging.warning(f'{self.name}: received {response.status_code} - {response.reason}')
            except requests.exceptions.ConnectionError as e:
                logging.warning(f'{self.name}: failed to connect to {self.address}. {e}')
            logging.warning(f'{self.name}: moving to next server')
            self.connecting = False
            self.switch()
        logging.warning(f'{self.name}: no working servers found')
        return None

    def report(self, output):
        logging.debug(f'Reporting {output}')
        for key, value in output.items():
            if key == 'transcoder_type_count':
                # we keep a list of all discovered modes so we can report zero when no session is running for a mode
                for mode in self.modes:
                    GAUGES[key].labels(self.name, mode).set(value[mode] if mode in value else 0)
            elif key == 'session_count':
                # we keep a list of all discovered users so we can report zero when a user is no longer logged in
                for user in self.users:
                    GAUGES[key].labels(self.name, user).set(value[user] if user in value else 0)
            elif key == 'version':
                GAUGES['server_info'].labels(value['server'], value['version']).set(1)
            else:
                GAUGES[key].labels(self.name).set(value)

    def process(self, output):
        logging.debug(f'Processing {output}')
        self.users.update(set([entry['user'] for entry in output['sessions']]))
        self.modes.update(set([entry['mode'] for entry in output['sessions'] if entry['transcode']]))
        return {
            'session_count': {
                user: len([entry for entry in output['sessions'] if entry['user'] == user])
                for user in self.users
            },
            'transcoder_count':
                len([entry for entry in output['sessions'] if entry['transcode']]),
            'transcoder_type_count': {
                mode: len([entry for entry in output['sessions'] if entry['mode'] == mode])
                for mode in self.modes},
            'transcoder_speed_total':
                sum([entry['speed'] for entry in output['sessions']]),
            'transcoder_encoding_count':
                len([entry for entry in output['sessions']
                     if entry['transcode'] and not entry['throttled']]),
            'version': output['version'],
        }

    @staticmethod
    def parse_session(session):
        if 'TranscodeSession' in session:
            return {
                'user': session['User']['title'],
                'transcode': True,
                'mode': session['TranscodeSession']['videoDecision'],
                'throttled': session['TranscodeSession']['throttled'],
                'speed': float(session['TranscodeSession']['speed'])
            }
        else:
            return {
                'user': session['User']['title'],
                'transcode': False,
                'mode': None,
                'throttled': False,
                'speed': 0
            }

    @staticmethod
    def parse_sessions(response):
        logging.debug(f'/status/session result: {response}')
        try:
            if response and 'Metadata' in response['MediaContainer']:
                return [PlexProbe.parse_session(session) for session in response['MediaContainer']['Metadata']]
        except KeyError as e:
            logging.warning(f'Failed to get sessions: missing {e}')
        except TypeError as e:
            logging.warning(f'Failed to get sessions: missing {e}')
        return []

    @staticmethod
    def parse_version(response):
        logging.debug(f'/identity result: {response}')
        if response:
            try:
                return response['MediaContainer']['version']
            except KeyError as e:
                logging.warning(f'Failed to get version: missing {e}')
        else:
            logging.warning(f'Failed to get version: no response received')
        return ''

    def measure_sessions(self):
        return PlexProbe.parse_sessions(self.call('/status/sessions'))

    def measure_version(self):
        return PlexProbe.parse_version(self.call('/identity'))

    def measure(self):
        return {
            'sessions': self.measure_sessions(),
            'version': {'server': 'plex', 'version': self.measure_version()},
        }


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
        servers = []
        if self.authtoken or self._login():
            headers = self.base_headers
            headers['X-Plex-Token'] = self.authtoken
            response = self.call('https://plex.tv/devices.xml', headers=headers)
            if response:
                servers = self._parse_servers(response)
        return servers

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
