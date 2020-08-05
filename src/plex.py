import requests
import logging
import xmltodict
import xml
from pimetrics.probe import APIProbe
import src.version
from prometheus_client import Gauge, start_http_server

GAUGES = {
    'plex_users': Gauge('mediaserver_plex_user_count', 'Active Plex viewing users', ['user']),
    'plex_clients': Gauge('mediaserver_plex_client_count', 'Active Plex viewing clients', ['client']),
    'plex_transcoder_count': Gauge('mediaserver_plex_transcoder_count', 'Active Transcoder count', ['mode']),
    'plex_transcoder_throttled_count':
        Gauge('mediaserver_plex_transcoder_throttled_count', 'Number of throttled transcoders', ['mode']),
    'plex_transcoder_speed': Gauge('mediaserver_plex_transcoder_speed', 'Speed of transcoders', ['mode']),
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
        self.address_index = (self.address_index+1) % len(self.addresses)

    @property
    def connecting(self):
        return self.status is None or self.status

    @connecting.setter
    def connecting(self, status):
        self.status = status


class PlexProbe(APIProbe, AddressManager):
    def __init__(self, authtoken, name, addresses, port=32400):
        AddressManager.__init__(self, addresses)
        APIProbe.__init__(self, '')
        self.name = name
        self.port = port
        self.headers = {
            'X-Plex-Token': authtoken,
            'Accept': 'application/json'
        }

    def call(self, endpoint):
        first_server = None
        while self.address != first_server:
            if first_server is None:
                first_server = self.address
            url = f'http://{self.address}:{self.port}{endpoint}'
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
        if output:
            logging.debug(output)
            for user, count in output['users'].items():
                GAUGES['plex_users'].labels(user).set(count)
            for client, count in output['clients'].items():
                GAUGES['plex_clients'].labels(client).set(count)
            for mode, attribs in output['transcoders'].items():
                GAUGES['plex_transcoder_count'].labels(mode).set(attribs['count'])
                GAUGES['plex_transcoder_throttled_count'].labels(mode).set(attribs['throttled'])
                GAUGES['plex_transcoder_speed'].labels(mode).set(attribs['speed'])

    def process(self, output):
        if output:
            logging.debug(output)
            users = set([entry['user'] for entry in output])
            clients = set([entry['client'] for entry in output])
            modes = set([entry['mode'] for entry in output if entry['transcode']])
            return {
                'users': {
                    user: len([entry for entry in output if entry['user'] == user]) for user in users
                },
                'clients': {
                    client: len([entry for entry in output if entry['client'] == client]) for client in clients
                },
                'transcoders': {
                    mode: {
                        'count': len([entry for entry in output if entry['mode']]),
                        'throttled': len([entry for entry in output if entry['mode'] == mode and entry['throttled']]),
                        'speed': sum([entry['speed'] for entry in output if entry['mode'] == mode])
                    } for mode in modes
                }
            }
        return {}

    @staticmethod
    def measure_session(session):
        values = {
            'client': session['Player']['product'],
            'user': session['User']['title'],
            'transcode': 'TranscodeSession' in session,
            'mode': None,
            'throttled': 0,
            'speed': 0}
        if 'TranscodeSession' in session:
            values['mode'] = session['TranscodeSession']['videoDecision']
            values['throttled'] = session['TranscodeSession']['throttled']
            values['speed'] = float(session['TranscodeSession']['speed'])
        return values

    def measure_sessions(self):
        try:
            response = self.call('/status/sessions')
            logging.debug(response)
            if 'Metadata' in response['MediaContainer']:
                return [self.measure_session(session) for session in response['MediaContainer']['Metadata']]
        except KeyError as e:
            logging.warning(f'Failed to get sessions: missing {e}')
        except TypeError as e:
            logging.warning(f'Failed to get sessions: missing {e}')
        return []

    def measure(self):
        return self.measure_sessions()


class PlexProbeManager:
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
    def _parse_servers(output, encoding):
        # output = '<?xml version="1.0" encoding="UTF-8"?>' \
        #          '<MediaContainer friendlyName="myPlex" identifier="com.plexapp.plugins.myplex" machineIdentifier="d16f8930d56547ff9ceabc1d09bd13bbec195c4b" size="2">' \
        #          '<Server accessToken="_o6iPBdHfAHM5n2YFv6G" name="Plex On Pi" address="81.164.72.79" port="80" version="1.19.5.3112-b23ab3896" scheme="http" host="81.164.72.79" localAddresses="192.168.0.10,10.0.6.34,172.18.0.8,10.0.8.22,10.0.0.54" machineIdentifier="9ee90fdb3204a90502f599a56f68eb07d54cc831" createdAt="1573850338" updatedAt="1596408667" owned="1" synced="0"/>' \
        #          '<Server accessToken="_o6iPBdHfAHM5n2YFv6G" name="Plex On Pi" address="81.164.72.79" port="80" version="1.19.5.3112-b23ab3896" scheme="http" host="81.164.72.79" localAddresses="192.168.0.10,10.0.6.34,172.18.0.8,10.0.8.22,10.0.0.54" machineIdentifier="9ee90fdb3204a90502f599a56f68eb07d54cc831" createdAt="1573850338" updatedAt="1596408667" owned="1" synced="0"/>' \
        #          '</MediaContainer>'
        # encoding = 'UTF-8'
        try:
            result = xmltodict.parse(output, encoding)
            size = int(result['MediaContainer']['@size'])
            servers = result['MediaContainer']['Server']

            if size == 1:
                return [{
                    'name': servers['@name'],
                    'addresses': servers['@localAddresses'].split(',')
                }]
            return [{
                'name': server['@name'],
                'addresses': server['@localAddresses'].split(',')
            } for server in servers]
        except KeyError as e:
            logging.warning(f'Failed to parse server list: missing tag {e}')
        except TypeError:
            logging.warning(f'Failed to parse server list. Unexpected tags found')
        except xml.parsers.expat.ExpatError as e:
            logging.warning(f'Failed to parse server list: {e}')
        return []

    def _get_servers(self):
        if not self.authtoken and not self._login():
            return []
        headers = self.base_headers
        headers['X-Plex-Token'] = self.authtoken
        response = requests.get('https://plex.tv/pms/servers.xml', headers=headers)
        if response.status_code == 200:
            logging.debug(response.content)
            return self._parse_servers(response.content, response.encoding)
        else:
            logging.error(f'Failed to retrieve server list from plex.tv: {response.status_code} - {response.reason}')
        return []

    def make_probes(self):
        self.probes = [PlexProbe(self.authtoken, server['name'], server['addresses']) for server in self._get_servers()]
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
