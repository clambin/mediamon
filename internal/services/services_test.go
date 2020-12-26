package services_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mediamon/internal/services"
	"os"
	"testing"
	"time"
)

func TestParsePartialConfig(t *testing.T) {
	var content = []byte(`
transmission:
  url: http://192.168.0.10:9091
  interval: '10s'
plex:
  username: email@example.com
  password: 'some-password'
  interval: '1m'
`)

	f, err := ioutil.TempFile("", "tmp")
	if err != nil {
		panic(err)
	}

	defer os.Remove(f.Name())
	_, _ = f.Write(content)
	_ = f.Close()

	cfg, err := services.ParseConfigFile(f.Name())

	assert.Nil(t, err)
	assert.Equal(t, "http://192.168.0.10:9091", cfg.Transmission.URL)
	assert.Equal(t, 10*time.Second, cfg.Transmission.Interval)
	assert.Equal(t, "", cfg.Sonarr.URL)
	assert.Equal(t, "", cfg.Sonarr.APIKey)
	assert.Equal(t, 30*time.Second, cfg.Sonarr.Interval)
	assert.Equal(t, "", cfg.Radarr.URL)
	assert.Equal(t, "", cfg.Radarr.APIKey)
	assert.Equal(t, 30*time.Second, cfg.Radarr.Interval)
	assert.Equal(t, "email@example.com", cfg.Plex.UserName)
	assert.Equal(t, "some-password", cfg.Plex.Password)
	assert.Equal(t, 1*time.Minute, cfg.Plex.Interval)
	assert.Equal(t, "", cfg.OpenVPN.Bandwidth.FileName)
	assert.Equal(t, 30*time.Second, cfg.OpenVPN.Bandwidth.Interval)
	assert.Equal(t, "", cfg.OpenVPN.Connectivity.Proxy)
	assert.Equal(t, "", cfg.OpenVPN.Connectivity.Token)
	assert.Equal(t, 5*time.Minute, cfg.OpenVPN.Connectivity.Interval)
}

func TestParseConfig(t *testing.T) {
	var content = []byte(`
transmission:
  url: http://192.168.0.10:9091
  interval: '10s'
sonarr:
  url: http://192.168.0.10:8989
  apikey: 'sonarr-api-key'
  interval: '5m'
radarr:
  url: http://192.168.0.10:7878
  apikey: 'radarr-api-key'
  interval: '5m'
plex:
  url: http://192.168.0.10:32400
  username: email@example.com
  password: 'some-password'
  interval: '1m'
openvpn:
  bandwidth:
    filename: /foo/bar
    interval: '30s'
  connectivity:
    proxy: http://localhost:8888
    token: 'some-token'
    interval: '5m'
futurefeature:
  foo: 'bar'
`)

	cfg, err := services.ParseConfig(content)

	assert.Nil(t, err)
	assert.Equal(t, "http://192.168.0.10:9091", cfg.Transmission.URL)
	assert.Equal(t, 10*time.Second, cfg.Transmission.Interval)
	assert.Equal(t, "http://192.168.0.10:8989", cfg.Sonarr.URL)
	assert.Equal(t, "sonarr-api-key", cfg.Sonarr.APIKey)
	assert.Equal(t, 5*time.Minute, cfg.Sonarr.Interval)
	assert.Equal(t, "http://192.168.0.10:7878", cfg.Radarr.URL)
	assert.Equal(t, "radarr-api-key", cfg.Radarr.APIKey)
	assert.Equal(t, 5*time.Minute, cfg.Radarr.Interval)
	assert.Equal(t, "http://192.168.0.10:32400", cfg.Plex.URL)
	assert.Equal(t, "email@example.com", cfg.Plex.UserName)
	assert.Equal(t, "some-password", cfg.Plex.Password)
	assert.Equal(t, 1*time.Minute, cfg.Plex.Interval)
	assert.Equal(t, "/foo/bar", cfg.OpenVPN.Bandwidth.FileName)
	assert.Equal(t, 30*time.Second, cfg.OpenVPN.Bandwidth.Interval)
	assert.Equal(t, "http://localhost:8888", cfg.OpenVPN.Connectivity.Proxy)
	assert.Equal(t, "some-token", cfg.OpenVPN.Connectivity.Token)
	assert.Equal(t, 5*time.Minute, cfg.OpenVPN.Connectivity.Interval)
}

func TestParseInvalidConfig(t *testing.T) {
	var content = []byte(`not a valid yaml file`)
	_, err := services.ParseConfig(content)
	assert.NotNil(t, err)
}

func TestParseMissingConfig(t *testing.T) {
	_, err := services.ParseConfigFile("not_a_file")
	assert.NotNil(t, err)
}
