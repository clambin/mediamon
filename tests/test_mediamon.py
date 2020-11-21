import argparse
from mediamon.mediamon import initialise, mediamon
from mediamon.configuration import get_configuration
from mediamon.transmission import TransmissionProbe
from mediamon.xxxarr import MonitorProbe
from mediamon.plex import PlexServer


def test_initialise():
    config = argparse.Namespace(interval=0, port=8080,
                                once=True, stub=True, debug=True,
                                services={
                                    'transmission': {
                                        'host': '192.168.0.10:9091',
                                        'interval': 1
                                    },
                                    'sonarr': {
                                        'host': '192.168.0.10:8989',
                                        'apikey': 'sonar-api-key',
                                        'interval': 2
                                    },
                                    'radarr': {
                                        'host': '192.168.0.10:7878',
                                        'apikey': 'radar-api-key',
                                        'interval': 3
                                    },
                                    'plex': {
                                        'username': 'email@example.com',
                                        'password': 'some-password',
                                        'interval': 4
                                    }
                                })

    scheduler = initialise(config)
    assert len(scheduler.scheduled_items) == 4
    assert type(scheduler.scheduled_items[0].probe) is TransmissionProbe
    assert scheduler.scheduled_items[0].interval == 1
    assert type(scheduler.scheduled_items[1].probe) is MonitorProbe
    assert scheduler.scheduled_items[1].probe.app == MonitorProbe.App.sonarr
    assert scheduler.scheduled_items[1].interval == 2
    assert type(scheduler.scheduled_items[2].probe) is MonitorProbe
    assert scheduler.scheduled_items[2].probe.app == MonitorProbe.App.radarr
    assert type(scheduler.scheduled_items[2].probe) is MonitorProbe
    assert scheduler.scheduled_items[2].interval == 3
    assert type(scheduler.scheduled_items[3].probe) is PlexServer
    assert scheduler.scheduled_items[3].interval == 4


def test_mediamon():
    configuration = get_configuration(['--once'])
    assert configuration.once is True
    assert mediamon(configuration) == 0
