import argparse
import yaml
import logging
import sys

from src.version import version


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

    # FIXME: allow credentials to be retrieved from file rather than cmdline options
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

    if config.services:
        try:
            with open(config.services, 'r') as f:
                config.services = yaml.safe_load(f)
        except FileNotFoundError as e:
            logging.critical(f'Could not open services file: {e}')
            sys.exit(1)
        except AttributeError as e:
            logging.critical(f'Could not parse services file: {e}')
            sys.exit(1)
    else:
        config.services = {}

    return config


def print_configuration(config):
    return ', '.join([f'{key}={val}' for key, val in vars(config).items()])
