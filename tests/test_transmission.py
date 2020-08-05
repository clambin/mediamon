from src.transmission import TransmissionProbe
from tests.utils import FakeResponse


class TransmissionTestProbe(TransmissionProbe):
    def __init__(self, host):
        super().__init__(host)

    def post(self, endpoint=None, headers=None, body=None):
        if 'X-Transmission-Session-Id' not in headers or headers['X-Transmission-Session-Id'] == '':
            return FakeResponse(409, {'X-Transmission-Session-Id': 'NewKey'}, {})
        else:
            return FakeResponse(200, {}, {'arguments': {
                "activeTorrentCount": 1,
                "cumulative-stats": {
                    "downloadedBytes": 259842832295,
                    "filesAdded": 218,
                    "secondsActive": 3106063,
                    "sessionCount": 19,
                    "uploadedBytes": 67534137454
                },
                "current-stats": {
                    "downloadedBytes": 53505238629,
                    "filesAdded": 39,
                    "secondsActive": 508967,
                    "sessionCount": 1,
                    "uploadedBytes": 14868574785
                },
                "downloadSpeed": 1000,
                "pausedTorrentCount": 2,
                "torrentCount": 3,
                "uploadSpeed": 500
            }})


def test_transmission():
    probe = TransmissionTestProbe('localhost:8080')
    probe.run()
    assert probe.measured()['activeTorrentCount'] == 1
    assert probe.measured()['pausedTorrentCount'] == 2
    assert probe.measured()['downloadSpeed'] == 1000
    assert probe.measured()['uploadSpeed'] == 500
