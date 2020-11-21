from mediamon.transmission import TransmissionProbe
from pimetrics.stubs import APIStub


class TransmissionTestProbe(APIStub, TransmissionProbe):
    def __init__(self, test_files):
        APIStub.__init__(self, test_files)
        TransmissionProbe.__init__(self, '')


testfiles = {
    'session-get': {
        'filename': 'samples/transmission-session-get.json',
    },
    'session-stats': {
        'filename': 'samples/transmission-session-stats.json',
    },
}


def test_transmission():
    probe = TransmissionTestProbe(testfiles)
    probe.run()
    measured = probe.measured()
    assert measured['active_torrent_count'] == 1
    assert measured['paused_torrent_count'] == 2
    assert measured['download_speed'] == 1000
    assert measured['upload_speed'] == 500
    assert measured['version'] == '2.94 (d8e60ee44f)'
