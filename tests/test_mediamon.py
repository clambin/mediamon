import argparse
from src.mediamon import mediamon, initialise
from src.mediacentre import TransmissionProbe, MonitorProbe


def test_initialise():
    config = argparse.Namespace(interval=0, port=8080,
                                transmission='http://localhost:8080',
                                sonarr='http://localhost:8081', sonarr_apikey='notreallyakey',
                                radarr='http://localhost:8082', radarr_apikey='notreallyakey',
                                once=True, stub=True, debug=True)
    scheduler = initialise(config)
    assert len(scheduler.scheduled_items) == 3
    assert type(scheduler.scheduled_items[0].probe) is TransmissionProbe
    assert type(scheduler.scheduled_items[1].probe) is MonitorProbe
    assert scheduler.scheduled_items[1].probe.app == MonitorProbe.App.sonarr
    assert type(scheduler.scheduled_items[2].probe) is MonitorProbe
    assert scheduler.scheduled_items[2].probe.app == MonitorProbe.App.radarr


def test_missing_sonarr_apikey():
    config = argparse.Namespace(interval=0, port=8080,
                                transmission='http://localhost:8080',
                                sonarr='http://localhost:8081', sonarr_apikey='',
                                radarr='http://localhost:8082', radarr_apikey='notreallyakey',
                                once=True, stub=True, debug=True)
    scheduler = initialise(config)
    assert len(scheduler.scheduled_items) == 2
    assert type(scheduler.scheduled_items[0].probe) is TransmissionProbe
    assert type(scheduler.scheduled_items[1].probe) is MonitorProbe
    assert scheduler.scheduled_items[1].probe.app == MonitorProbe.App.radarr


def test_missing_radarr_apikey():
    config = argparse.Namespace(interval=0, port=8080,
                                transmission='http://localhost:8080',
                                sonarr='http://localhost:8081', sonarr_apikey='notreallyakey',
                                radarr='http://localhost:8082', radarr_apikey='',
                                once=True, stub=True, debug=True)
    scheduler = initialise(config)
    assert len(scheduler.scheduled_items) == 2
    assert type(scheduler.scheduled_items[0].probe) is TransmissionProbe
    assert type(scheduler.scheduled_items[1].probe) is MonitorProbe
    assert scheduler.scheduled_items[1].probe.app == MonitorProbe.App.sonarr


# def test_mediamon():
#    config = argparse.Namespace(interval=0, port=8080,
#                                transmission='http://localhost:8080',
#                                sonarr='http://localhost:8081', sonarr_apikey='notreallyakey',
#                                radarr='http://localhost:8082', radarr_apikey='notreallyakey',
#                                once=True, stub=True, debug=True)
#    assert mediamon(config) == 0
