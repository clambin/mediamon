import argparse
from mediamon.mediamon import initialise
from mediamon.transmission import TransmissionProbe
from mediamon.xxxarr import MonitorProbe
from mediamon.plex import PlexServer


def test_initialise():
    config = argparse.Namespace(interval=0, port=8080,
                                once=True, stub=True, debug=True,
                                services={
                                    'transmission': {'host': '192.168.0.10:9091'},
                                    'sonarr': {'host': '192.168.0.10:8989', 'apikey': 'sonar-api-key'},
                                    'radarr': {'host': '192.168.0.10:7878', 'apikey': 'radar-api-key'},
                                    'plex': {'username': 'email@example.com', 'password': 'some-password'}
                                })

    scheduler = initialise(config)
    assert len(scheduler.scheduled_items) == 4
    assert type(scheduler.scheduled_items[0].probe) is TransmissionProbe
    assert type(scheduler.scheduled_items[1].probe) is MonitorProbe
    assert scheduler.scheduled_items[1].probe.app == MonitorProbe.App.sonarr
    assert type(scheduler.scheduled_items[2].probe) is MonitorProbe
    assert scheduler.scheduled_items[2].probe.app == MonitorProbe.App.radarr
    assert type(scheduler.scheduled_items[2].probe) is MonitorProbe
    assert type(scheduler.scheduled_items[3].probe) is PlexServer
