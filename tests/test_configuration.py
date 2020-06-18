import argparse
import pytest
from src.configuration import str2bool, get_configuration


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
    assert config.transmission == ''
    assert config.sonarr == ''
    assert config.sonarr_apikey == ''
    assert config.radarr == ''
    assert config.radarr_apikey == ''

