import logging
from prometheus_client import start_http_server
from src.version import version
from src.configuration import print_configuration
from src.mediacentre import TransmissionProbe, MonitorProbe
from src.plex import PlexServer
from pimetrics.scheduler import Scheduler


def initialise(config):
    scheduler = Scheduler()

    if config.transmission:
        scheduler.register(TransmissionProbe(config.transmission), 5)
    if config.sonarr:
        if config.sonarr_apikey:
            scheduler.register(MonitorProbe(config.sonarr, MonitorProbe.App.sonarr, config.sonarr_apikey), 300)
        else:
            logging.warning('sonarr url specified but apikey missing. Ignoring')
    if config.radarr:
        if config.radarr_apikey:
            scheduler.register(MonitorProbe(config.radarr, MonitorProbe.App.radarr, config.radarr_apikey), 300)
        else:
            logging.warning('radarr url specified but apikey missing. Ignoring')
    if config.plex_username and config.plex_password:
        scheduler.register(PlexServer(config.plex_username, config.plex_password))
    return scheduler


def mediamon(config):
    logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s', datefmt='%Y-%m-%d %H:%M:%S',
                        level=logging.DEBUG if config.debug else logging.INFO)
    logging.info(f'Starting mediamon v{version}')
    logging.info(f'Configuration: {print_configuration(config)}')

    start_http_server(config.port)

    scheduler = initialise(config)
    if config.once:
        scheduler.run(once=True)
    else:
        while True:
            scheduler.run(duration=config.interval)
    return 0
