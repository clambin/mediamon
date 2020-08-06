import logging
import prometheus_client
from src.version import version
from src.configuration import print_configuration
from src.transmission import TransmissionProbe
from src.xxxarr import MonitorProbe
from src.plex import PlexServer
from pimetrics.scheduler import Scheduler


def initialise(config):
    scheduler = Scheduler()

    if 'transmission' in config.services:
        try:
            scheduler.register(
                TransmissionProbe(config.services['transmission']['host']),
                5
            )
        except KeyError as e:
            logging.warning(f'transmission config missing {e}. Skipping')

    if 'sonarr' in config.services:
        try:
            scheduler.register(
                MonitorProbe(
                    config.services['sonarr']['host'],
                    MonitorProbe.App.sonarr,
                    config.services['sonarr']['apikey']
                ),
                300
            )
        except KeyError as e:
            logging.warning(f'sonarr config missing {e}. Skipping')

    if 'radarr' in config.services:
        try:
            scheduler.register(
                MonitorProbe(
                    config.services['radarr']['host'],
                    MonitorProbe.App.radarr,
                    config.services['radarr']['apikey']),
                300
            )
        except KeyError as e:
            logging.warning(f'radarr config missing {e}. Skipping')

    if 'plex' in config.services:
        try:
            scheduler.register(
                PlexServer(config.services['plex']['username'], config.services['plex']['password']),
                5
            )
        except KeyError as e:
            logging.warning(f'plex config missing {e}. Skipping')

    if len(scheduler.scheduled_items) == 0:
        logging.error('No services defined')
        # TODO: fix this in pimetrics
        scheduler.min_interval = 5

    return scheduler


def mediamon(config):
    logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s', datefmt='%Y-%m-%d %H:%M:%S',
                        level=logging.DEBUG if config.debug else logging.INFO)
    logging.info(f'Starting mediamon v{version}')
    logging.info(f'Configuration: {print_configuration(config)}')

    prometheus_client.start_http_server(config.port)

    scheduler = initialise(config)
    if config.once:
        scheduler.run(once=True)
    else:
        while True:
            scheduler.run(duration=config.interval)
    return 0
