import json
from mediamon.xxxarr import MonitorProbe
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


# TODO: align plex/xxxarr stubbing approaches
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
    assert measured['xxxarr_calendar'] == 1
    assert measured['xxxarr_queue'] == 1
    assert measured['xxxarr_monitored'] == 1
    assert measured['xxxarr_unmonitored'] == 1
    assert measured['version'] == '2.0.0.5344'


def test_radarr():
    probe = MonitorTestProbe('', MonitorProbe.App.radarr, '', radarr_test_files)
    assert probe.name == 'radarr'
    probe.connecting = False
    probe.run()
    assert probe.connecting is True
    measured = probe.measured()
    assert measured
    assert measured['xxxarr_calendar'] == 0
    assert measured['xxxarr_queue'] == 2
    assert measured['xxxarr_monitored'] == 1
    assert measured['xxxarr_unmonitored'] == 1
    assert measured['version'] == '0.2.0.1504'
