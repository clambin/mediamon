import json
from src.xxxarr import MonitorProbe
from tests.utils import FakeResponse

sonarr_test_files = {
    'api/system/status': 'samples/sonarr_version.json',
    'api/calendar': 'samples/sonarr_calendar.json',
    'api/queue': 'samples/sonarr_queue.json',
    'api/series': 'samples/sonarr_series.json',
}

radarr_test_files = {
    'api/system/status': 'samples/radarr_version.json',
    'api/calendar': 'samples/radarr_calendar.json',
    'api/queue': 'samples/radarr_queue.json',
    'api/movie': 'samples/radarr_movie.json',
}


class MonitorTestProbe(MonitorProbe):
    def __init__(self, host, name, api_key, testfiles=None):
        self.testfiles = testfiles if testfiles is not None else dict()
        super().__init__(host, name, api_key)

    def get(self, endpoint=None, headers=None, body=None, params=None):
        if endpoint in self.testfiles:
            with open(self.testfiles[endpoint], 'r') as f:
                return FakeResponse(200, {}, json.load(f))
        else:
            return FakeResponse(404, {}, '')


def test_sonarr():
    probe = MonitorTestProbe('', MonitorProbe.App.sonarr, '', sonarr_test_files)
    assert probe.name == 'sonarr'
    probe.connecting = False
    probe.run()
    assert probe.connecting is True
    measured = probe.measured()
    assert measured
    assert measured['calendar'] == 1
    assert measured['queue'] == 1
    assert measured['monitored'] == (1, 1)
    assert measured['version'] == '2.0.0.5344'


def test_radarr():
    probe = MonitorTestProbe('', MonitorProbe.App.radarr, '', radarr_test_files)
    assert probe.name == 'radarr'
    probe.connecting = False
    probe.run()
    assert probe.connecting is True
    measured = probe.measured()
    assert measured
    assert measured['calendar'] == 0
    assert measured['queue'] == 2
    assert measured['monitored'] == (1, 1)
    assert measured['version'] == '0.2.0.1504'
