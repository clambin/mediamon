import argparse
import pytest
from mediamon.configuration import str2bool, get_configuration, print_configuration


def test_str2bool():
    assert str2bool(True) is True
    for arg in ['yes', 'true', 't', 'y', '1', 'on']:
        assert str2bool(arg) is True
    for arg in ['no', 'false', 'f', 'n', '0', 'off']:
        assert str2bool(arg) is False
    with pytest.raises(argparse.ArgumentTypeError) as e:
        assert str2bool('maybe')
    assert str(e.value) == 'Boolean value expected.'


def test_main_config():
    args = '--interval 25 --port 1234 --once --debug'.split()
    config = get_configuration(args)
    assert config.interval == 25
    assert config.port == 1234
    assert config.once
    assert config.debug


def test_default_config():
    config = get_configuration([])
    assert config.debug is False
    assert config.interval == 5
    assert config.port == 8080
    assert config.stub is False
    assert config.services == {}


def test_print_config():
    args = '--file samples/services.yml'.split()
    config = get_configuration(args)
    assert config.services is not None
    output = print_configuration(config)
    assert output == "interval=5, port=8080, once=False, stub=False, debug=False, " \
                     "file=samples/services.yml, " \
                     "services={" \
                     "'transmission': {'host': '192.168.0.10:9091', " \
                     "'interval': 10}, " \
                     "'sonarr': {'host': '192.168.0.10:8989', 'apikey': '********************************', " \
                     "'interval': 300}, " \
                     "'radarr': {'host': '192.168.0.10:7878', 'apikey': '********************************', " \
                     "'interval': 500}, " \
                     "'plex': {'username': 'email@example.com', 'password': '************', " \
                     "'interval': 5}}"


def test_services():
    args = '--file samples/services.yml'.split()
    config = get_configuration(args)
    assert config.services['transmission']['host'] == '192.168.0.10:9091'
    assert config.services['transmission']['interval'] == 10
    assert config.services['sonarr']['host'] == '192.168.0.10:8989'
    assert config.services['sonarr']['apikey'] == 'sonar-api-key'
    assert config.services['sonarr']['interval'] == 300
    assert config.services['radarr']['host'] == '192.168.0.10:7878'
    assert config.services['radarr']['apikey'] == 'radar-api-key'
    assert config.services['radarr']['interval'] == 500
    assert config.services['plex']['username'] == 'email@example.com'
    assert config.services['plex']['password'] == 'some-password'
    assert config.services['plex']['interval'] == 5
    alt_args = '-f samples/services.yml'.split()
    alt_config = get_configuration(alt_args)
    assert config == alt_config


def test_invalid_services():
    args = '-f samples/no_services.yml'.split()
    config = get_configuration(args)
    assert config.services == {}
    args = '-f samples/invalid-services.yml'.split()
    config = get_configuration(args)
    assert config.services == {}
