import argparse

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
    # Media server monitoring
    parser.add_argument('--transmission', default='',
                        help='Transmission address (<host>:<port>)')
    parser.add_argument('--sonarr', default='',
                        help='Sonarr address (<host>:<port>)')
    parser.add_argument('--sonarr-apikey', default='',
                        help='Sonarr API Key')
    parser.add_argument('--radarr', default='',
                        help='Radarr address (<host>:<port>)')
    parser.add_argument('--radarr-apikey', default='',
                        help='Radarr API Key')
    parser.add_argument('--plex-username', default='',
                        help='Plex username')
    parser.add_argument('--plex-password', default='',
                        help='Plex password')
    return parser.parse_args(args)


def print_configuration(config):
    return ', '.join([f'{key}={val}' for key, val in vars(config).items()])
