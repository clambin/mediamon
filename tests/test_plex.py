import json

from mediamon.plex import PlexProbe, AddressManager, PlexServer
from pimetrics.stubs import APIStub


def test_address_manager():
    mgr = AddressManager(['a', 'b', 'c'])
    assert mgr.healthy is True
    assert mgr.address == 'a'
    mgr.switch()
    assert mgr.address == 'b'
    mgr.switch()
    assert mgr.address == 'c'
    mgr.switch()
    assert mgr.address == 'a'
    mgr.healthy = False
    assert mgr.healthy is False


def test_parse_sessions():
    def get_sessions(filename):
        with open(filename, 'r') as f:
            return PlexProbe.parse_sessions(json.loads(f.read()))
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


class PlexTestProbe(APIStub, PlexProbe):
    def __init__(self, authtoken, name, addresses, testfiles=None):
        PlexProbe.__init__(self, authtoken, name, addresses)
        APIStub.__init__(self, testfiles)


plex_responses = {
    '/status/sessions': {
        'filename': 'samples/plex_session_multiple.json',
    },
    '/identity': {
        'filename': 'samples/plex_identity.json',
    },
}


def test_plexprobe():
    probe = PlexTestProbe('', 'Plex', ['foo', 'bar'], plex_responses)
    probe.run()
    measured = probe.measured()
    assert measured == {
        'plex_session_count': {'foo': 1, 'bar': 1},
        'plex_transcoder_count': 1,
        'plex_transcoder_encoding_count': 1,
        'plex_transcoder_speed_total': 3.1,
        'plex_transcoder_type_count': {'copy': 1},
        'version': '1.20.2.3402-0fec14d92'
    }


plex_server_responses = {
    '/users/sign_in.xml': {
        'filename': 'samples/plex_sign_in.xml',
        'raw': True,
    },
    '/devices.xml': {
        'filename': 'samples/plex_devices.xml',
        'raw': True,
    }
}


class PlexTestServer(APIStub, PlexServer):
    def __init__(self, username, password, testfiles=None):
        APIStub.__init__(self, testfiles)
        PlexServer.__init__(self, username, password)


def test_plexserver():
    server = PlexTestServer('', '', plex_server_responses)
    server.run()
    assert server.authtoken == '!!!!!!!!!!!!!!!!!!!!'
    assert len(server.probes) == 2
    assert set([probe.name for probe in server.probes]) == {'Plex Server 1', 'Plex Server 2'}
    assert server.probes[0].name == 'Plex Server 1'
    assert server.probes[0].addresses == [
        'http://127.0.0.1:32400',
        'http://127.0.0.1:32401'
    ]
    assert server.probes[1].name == 'Plex Server 2'
    assert server.probes[1].addresses == [
        'http://127.0.0.1:32402'
    ]
    oldprobes = set(server.probes)
    server.probes[1].healthy = False
    server._healthcheck()
    assert len(server.probes) == 2
    newprobes = set(server.probes)
    assert oldprobes != newprobes
