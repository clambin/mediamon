import argparse
import logging
import copy

import yaml
import yaml.scanner

from mediamon.version import version


def str2bool(v):
    if isinstance(v, bool):
        return v
    if v.lower() in ('yes', 'true', 't', 'y', '1', 'on'):
        return True
    elif v.lower() in ('no', 'false', 'f', 'n', '0', 'off'):
        return False
    else:
        raise argparse.ArgumentTypeError('Boolean value expected.')


def get_configuration(args=None):
    default_interval = 5
    default_port = 8080

    parser = argparse.ArgumentParser()
    parser.add_argument('--version', action='version', version=f'%(prog)s {version}')
    parser.add_argument('--interval', type=int, default=default_interval,
                        help=f'Time between measurements (default: {default_interval} sec)')
    parser.add_argument('--port', type=int, default=default_port,
                        help=f'Prometheus listener port (default: {default_port})')
    parser.add_argument('--once', action='store_true',
                        help='Measure once and then terminate')
    parser.add_argument('--stub', action='store_true',
                        help='Use stubs (for debugging only')
    parser.add_argument('--debug', action='store_true',
                        help='Set logging level to debug')
    parser.add_argument('--services', default='',
                        help='Service configuration file')
    config = parser.parse_args(args)

    services_filename = config.services
    config.services = {}
    if services_filename:
        try:
            with open(services_filename, 'r') as f:
                config.services = yaml.safe_load(f)
        except FileNotFoundError as e:
            logging.critical(f'Could not open services file: {e}')
        except yaml.scanner.ScannerError as e:
            logging.critical(f'Could not parse services file: {e}')

    return config


def print_configuration(config):
    redacted = copy.deepcopy(config)
    if redacted.services:
        if 'sonarr' in redacted.services:
            redacted.services['sonarr']['apikey'] = '*' * 32
        if 'radarr' in redacted.services:
            redacted.services['radarr']['apikey'] = '*' * 32
        if 'plex' in redacted.services:
            redacted.services['plex']['password'] = '*' * 12
    return ', '.join([f'{key}={val}' for key, val in vars(redacted).items()])
