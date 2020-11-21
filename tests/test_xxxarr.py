from mediamon.xxxarr import MonitorProbe
from pimetrics.stubs import APIStub


sonarr_test_files = {
    'api/system/status': {'filename': 'samples/sonarr_version.json'},
    'api/calendar': {'filename': 'samples/sonarr_calendar.json'},
    'api/queue': {'filename': 'samples/sonarr_queue.json'},
    'api/series': {'filename': 'samples/sonarr_series.json'},
}

radarr_test_files = {
    'api/system/status': {'filename': 'samples/radarr_version.json'},
    'api/calendar': {'filename': 'samples/radarr_calendar.json'},
    'api/queue': {'filename': 'samples/radarr_queue.json'},
    'api/movie': {'filename': 'samples/radarr_movie.json'},
}


class MonitorTestProbe(APIStub, MonitorProbe):
    def __init__(self, host, name, api_key, testfiles=None):
        APIStub.__init__(self, testfiles)
        MonitorProbe.__init__(self, host, name, api_key)


def test_sonarr():
    probe = MonitorTestProbe('', MonitorProbe.App.sonarr, '', sonarr_test_files)
    assert probe.name == 'sonarr'
    probe.healthy = False
    probe.run()
    assert probe.healthy is True
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
    probe.healthy = False
    probe.run()
    assert probe.healthy is True
    measured = probe.measured()
    assert measured
    assert measured['xxxarr_calendar'] == 0
    assert measured['xxxarr_queue'] == 2
    assert measured['xxxarr_monitored'] == 1
    assert measured['xxxarr_unmonitored'] == 1
    assert measured['version'] == '0.2.0.1504'
