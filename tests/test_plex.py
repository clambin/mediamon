import json
import logging

from src.plex import PlexProbe, AddressManager, PlexServer
from tests.utils import APIStub


class PlexTestProbe(APIStub, PlexProbe):
    def __init__(self, authtoken, name, addresses):
        PlexProbe.__init__(self, authtoken, name, addresses)
        APIStub.__init__(self)


class PlexServerTest(APIStub, PlexServer):
    def __init__(self, username, password):
        PlexServer.__init__(self, username, password)
        APIStub.__init__(self)

    def _login(self):
        return True


def test_address_manager():
    mgr = AddressManager(['a', 'b', 'c'])
    assert mgr.connecting is True
    assert mgr.address == 'a'
    mgr.switch()
    assert mgr.address == 'b'
    mgr.switch()
    assert mgr.address == 'c'
    mgr.switch()
    assert mgr.address == 'a'
    mgr.connecting = False
    assert mgr.connecting is False


def get_sessions(filename):
    with open(filename, 'r') as f:
        content = json.loads(f.read())
        return PlexProbe.parse_sessions(content)


def test_parse_sessions():
    logging.basicConfig(level=logging.DEBUG)
    sessions = get_sessions('samples/plex_session_empty.json')
    assert len(sessions) == 0
    sessions = get_sessions('samples/plex_session_player.json')
    assert len(sessions) == 1
    assert sessions[0]['transcode'] is False
    assert sessions[0]['mode'] is None
    assert sessions[0]['throttled'] is False
    assert sessions[0]['speed'] == 0
    sessions = get_sessions('samples/plex_session_web.json')
    assert len(sessions) == 1
    assert sessions[0]['transcode'] is True
    assert sessions[0]['mode'] == 'copy'
    assert sessions[0]['throttled'] is False
    assert sessions[0]['speed'] == 3.1
    sessions = get_sessions('samples/plex_session_multiple.json')
    assert len(sessions) == 2
    assert sessions[0]['transcode'] is False
    assert sessions[0]['mode'] is None
    assert sessions[0]['throttled'] is False
    assert sessions[0]['speed'] == 0
    assert sessions[1]['transcode'] is True
    assert sessions[1]['mode'] == 'copy'
    assert sessions[1]['throttled'] is False
    assert sessions[1]['speed'] == 3.1


def test_plexprobe_parser():
    probe = PlexProbe(None, None, None)
    sessions = get_sessions('samples/plex_session_multiple.json')
    output = probe.process(sessions)
    assert output['session_count'] == 2
    assert output['transcoder_count'] == 1
    assert output['transcoder_type_count'] == {'copy': 1}
    assert output['transcoder_speed_total'] == 3.1
    assert output['transcoder_encoding_count'] == 1


def test_plexserver_parser():
    with open('samples/plex_devices.xml') as f:
        content = f.read()
        servers = PlexServer._parse_servers(content, 'UTF-8')
        assert len(servers) == 2
        assert servers[0] == {'name': 'Plex Server 1', 'addresses': [
            'http://192.168.0.10:32400',
            'http://10.0.6.34:32400',
            'http://172.18.0.8:32400',
            'http://10.0.8.22:32400',
            'http://10.0.0.54:32400'
        ]}
        assert servers[1] == {'name': 'Plex Server 2', 'addresses': [
            'http://172.18.0.9:32400'
        ]}


def test_plexprobe():
    probe = PlexTestProbe('', '', '')
    probe.set_response('samples/plex_session_multiple.json', mode=APIStub.Mode.json)
    probe.run()
    measured = probe.measured()
    assert measured == {
        'session_count': 2,
        'transcoder_count': 1,
        'transcoder_encoding_count': 1,
        'transcoder_speed_total': 3.1,
        'transcoder_type_count': {'copy': 1}
    }


def test_plexserver():
    server = PlexServerTest('', '')
    server.set_response('samples/plex_devices.xml')
    probes = server.make_probes()
    assert len(probes) == 2
    oldprobes = set(probes)
    probes[1].connecting = False
    server._healthcheck()
    probes = server.probes
    assert len(probes) == 2
    newprobes = set(probes)
    assert oldprobes != newprobes
    assert set([probe.name for probe in probes]) == {'Plex Server 1', 'Plex Server 2'}
