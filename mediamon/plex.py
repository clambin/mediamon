import requests
import logging
import xmltodict
import xml
from collections import OrderedDict
from pimetrics.probe import APIProbe
import mediamon.version
from mediamon import metrics


class AddressManager:
    def __init__(self, addresses):
        self.addresses = addresses
        self.address_index = 0
        self.healthy = True

    @property
    def address(self):
        if self.addresses and len(self.addresses) > 0:
            return self.addresses[self.address_index]
        return None

    def switch(self):
        self.address_index = (self.address_index + 1) % len(self.addresses)


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

    def apicall(self, endpoint):
        first_server = None
        while self.address != first_server:
            try:
                if first_server is None:
                    first_server = self.address
                self.url = self.address
                response = self.call(endpoint, headers=self.headers)
                if response:
                    if self.healthy is False:
                        logging.info(f'{self.name}: connection established on {self.address}')
                        self.healthy = True
                    return response
                else:
                    logging.warning(f'{self.name}: failed to get data from {self.address}')
            except requests.exceptions.ConnectionError as e:
                logging.warning(f'{self.name}: failed to connect to {self.address}. {e}')
            logging.warning(f'{self.name}: moving to next server')
            self.healthy = False
            self.switch()
        logging.warning(f'{self.name}: no working servers found')
        return None

    def report(self, output):
        logging.debug(f'Reporting {output}')
        # we keep a list of all discovered modes so we can report zero when no session is running for a mode
        for mode in self.modes:
            if mode not in output['plex_transcoder_type_count']:
                output['plex_transcoder_type_count'][mode] = 0
        # we keep a list of all discovered users so we can report zero when a user is no longer logged in
        for user in self.users:
            if user not in output['plex_session_count']:
                output['plex_session_count'][user] = 0
        metrics.report(output, 'plex')

    def process(self, output):
        logging.debug(f'Processing {output}')
        self.users.update(set([entry['user'] for entry in output['sessions']]))
        self.modes.update(set([entry['mode'] for entry in output['sessions'] if entry['transcode']]))
        return {
            'plex_session_count': {
                user: len([entry for entry in output['sessions'] if entry['user'] == user])
                for user in self.users
            },
            'plex_transcoder_count':
                len([entry for entry in output['sessions'] if entry['transcode']]),
            'plex_transcoder_type_count': {
                mode: len([entry for entry in output['sessions'] if entry['mode'] == mode])
                for mode in self.modes},
            'plex_transcoder_speed_total':
                sum([entry['speed'] for entry in output['sessions']]),
            'plex_transcoder_encoding_count':
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
            logging.warning('Failed to get version: no response received')
        return ''

    def measure_sessions(self):
        return PlexProbe.parse_sessions(self.apicall('/status/sessions'))

    def measure_version(self):
        return PlexProbe.parse_version(self.apicall('/identity'))

    def measure(self):
        return {
            'sessions': self.measure_sessions(),
            'version': self.measure_version(),
        }


class PlexServer(APIProbe):
    def __init__(self, username, password):
        super().__init__('https://plex.tv')
        self.authtext = f'user%5Blogin%5D={username}&user%5Bpassword%5D={password}'
        self.base_headers = {
            'X-Plex-Product': 'mediamon',
            'X-Plex-Version': mediamon.version.version,
            # FIXME: generate UUID
            'X-Plex-Client-Identifier': f'mediamon-v{mediamon.version.version}'
        }
        self.authtoken = None
        self.probes = []

    def measure(self):
        raise AssertionError('should never be called')

    def call(self, endpoint=None, headers=None, body=None, params=None, method=APIProbe.Method.GET):
        try:
            if method == APIProbe.Method.GET:
                response = self.get(endpoint=endpoint, headers=headers, body=body, params=params)
                if response.status_code == 200:
                    return response.content
            else:
                # FIXME: extend APIProbe to use data rather than json to post content?
                response = requests.post(f'{self.url}{endpoint}', headers=headers, data=body, params=params)
                if response.status_code == 201:
                    return response.content
            logging.error("%d - %s" % (response.status_code, response.reason))
        except requests.exceptions.RequestException as err:
            logging.warning(f'Failed to call "{self.url}": "{err}')
        return None

    def apicall(self, endpoint, headers=None, body=None, method=APIProbe.Method.GET):
        try:
            return self.call(endpoint, headers=headers, body=body, method=method)
        except requests.exceptions.ConnectionError as e:
            logging.warning(f'Failed to connect to {self.url}{endpoint}: {e}')
        return None

    def _login(self):
        try:
            content = self.apicall('/users/sign_in.xml', headers=self.base_headers, body=self.authtext,
                                   method=APIProbe.Method.POST)
            if content:
                try:
                    result = xmltodict.parse(content, 'UTF-8')
                    self.authtoken = result['user']['@authenticationToken']
                    return True
                except KeyError as e:
                    logging.error(f'Could not parse login response: {e}')
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

    def _get_servers(self):
        servers = []
        if self.authtoken or self._login():
            headers = self.base_headers
            headers['X-Plex-Token'] = self.authtoken
            content = self.apicall('/devices.xml', headers=headers)
            if content:
                servers = self._parse_servers(content)
        return servers

    def _make_probes(self):
        servers = self._get_servers()
        for server in servers:
            logging.info(f'Plex server found: {server["name"]}: {",".join(server["addresses"])}')
        self.probes = [PlexProbe(self.authtoken, server['name'], server['addresses']) for server in servers]
        return self.probes

    def _healthcheck(self):
        if unhealthy := [probe for probe in self.probes if probe.healthy is False]:
            healthy = [probe for probe in self.probes if probe.healthy]
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
            self._make_probes()
        for probe in self.probes:
            probe.run()
        self._healthcheck()
