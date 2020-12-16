package services_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mediamon/internal/services"
	"os"
	"testing"
)

func TestParsePartialConfig(t *testing.T) {
	var content = []byte(`
transmission:
  url: http://192.168.0.10:9091
  interval: 10s
plex:
  username: email@example.com
  password: 'some-password'
  interval: 1m
`)

	f, err := ioutil.TempFile("", "tmp")
	if err != nil {
		panic(err)
	}

	defer os.Remove(f.Name())
	_, _ = f.Write(content)
	_ = f.Close()

	var cfg = services.Config{}

	err = services.ParseConfigFile(f.Name(), &cfg)

	assert.Nil(t, err)
	assert.Equal(t, "http://192.168.0.10:9091", cfg.Transmission.URL)
	assert.Equal(t, "10s", cfg.Transmission.Interval)
	assert.Equal(t, "", cfg.Sonarr.URL)
	assert.Equal(t, "", cfg.Sonarr.APIKey)
	assert.Equal(t, "", cfg.Sonarr.Interval)
	assert.Equal(t, "", cfg.Radarr.URL)
	assert.Equal(t, "", cfg.Radarr.APIKey)
	assert.Equal(t, "", cfg.Radarr.Interval)
	assert.Equal(t, "email@example.com", cfg.Plex.UserName)
	assert.Equal(t, "some-password", cfg.Plex.Password)
	assert.Equal(t, "1m", cfg.Plex.Interval)
}
func TestParseConfig(t *testing.T) {
	var content = []byte(`
transmission:
  url: http://192.168.0.10:9091
  interval: 10s
sonarr:
  url: http://192.168.0.10:8989
  apikey: 'sonarr-api-key'
  interval: 5m
radarr:
  url: http://192.168.0.10:7878
  apikey: 'radarr-api-key'
  interval: 5m
plex:
  url: http://192.168.0.10:32400
  username: email@example.com
  password: 'some-password'
  interval: 1m
`)

	var cfg = services.Config{}

	err := services.ParseConfig(content, &cfg)

	assert.Nil(t, err)
	assert.Equal(t, "http://192.168.0.10:9091", cfg.Transmission.URL)
	assert.Equal(t, "10s", cfg.Transmission.Interval)
	assert.Equal(t, "http://192.168.0.10:8989", cfg.Sonarr.URL)
	assert.Equal(t, "sonarr-api-key", cfg.Sonarr.APIKey)
	assert.Equal(t, "5m", cfg.Sonarr.Interval)
	assert.Equal(t, "http://192.168.0.10:7878", cfg.Radarr.URL)
	assert.Equal(t, "radarr-api-key", cfg.Radarr.APIKey)
	assert.Equal(t, "5m", cfg.Radarr.Interval)
	assert.Equal(t, "http://192.168.0.10:32400", cfg.Plex.URL)
	assert.Equal(t, "email@example.com", cfg.Plex.UserName)
	assert.Equal(t, "some-password", cfg.Plex.Password)
	assert.Equal(t, "1m", cfg.Plex.Interval)
}

func TestParseInvalidConfig(t *testing.T) {
	var content = []byte(`not a valid yaml file`)
	var cfg = services.Config{}

	err := services.ParseConfig(content, &cfg)

	assert.NotNil(t, err)
}
