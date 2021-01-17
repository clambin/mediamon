package mediamon_test

import (
	"github.com/clambin/mediamon/internal/mediamon"
	"github.com/clambin/mediamon/internal/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartProbes(t *testing.T) {
	const servicesYaml = `transmission:
  url: http://127.0.0.1:8000/transmission/rpc
  interval: 5m
sonarr:
  url: http://127.0.0.1:8001
  interval: 5m
  apikey: somekey
radarr:
  url: http://127.0.0.1:8002
  interval: 5m
  apikey: somekey
plex:
  url: http://127.0.0.1:8003
  interval: 5m
  username: user@example.com
  password: somepassword
openvpn:
  bandwidth:
    filename: somefile
    interval: 5m
  connectivity:
    proxy: http://127.0.0.1:8004
    token: 'sometoken'
    interval: 5m`

	svcs, err := services.ParseConfig([]byte(servicesYaml))

	assert.Nil(t, err)

	cfg := mediamon.Configuration{
		Port:     0,
		Debug:    false,
		Services: svcs,
	}

	expected := []string{
		"Transmission", "Sonarr", "Radarr", "Plex", "Bandwidth", "Connectivity",
	}

	probes := mediamon.StartProbes(&cfg)
	assert.Len(t, probes, len(expected))
	for _, val := range expected {
		assert.Contains(t, probes, val)
	}
}
