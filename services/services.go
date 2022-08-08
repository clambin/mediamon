package services

import (
	"fmt"
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"time"
)

// Config contains the different possible services for mediamon to monitor
type Config struct {
	Transmission transmission.Config
	Sonarr       xxxarr.Config
	Radarr       xxxarr.Config
	Plex         plex.Config
	OpenVPN      struct {
		Bandwidth    bandwidth.Config
		Connectivity connectivity.Config
	}
}

// ParseConfigFile reads the configuration from the specified yaml file
func ParseConfigFile(fileName string) (config *Config, err error) {
	var cfg []byte
	if cfg, err = os.ReadFile(fileName); err != nil {
		return
	}
	cfg = []byte(os.ExpandEnv(string(cfg)))

	config = &Config{}
	config.OpenVPN.Connectivity.Interval = 5 * time.Minute
	err = yaml.Unmarshal(cfg, config)

	// check proxy URL here, so we can raise an error before trying to start the collectors
	if config.OpenVPN.Connectivity.Proxy != "" {
		var proxy *url.URL
		if proxy, err = url.Parse(config.OpenVPN.Connectivity.Proxy); err != nil {
			err = fmt.Errorf("invalid proxy url: %w", err)
		} else if proxy.Scheme == "" || proxy.Host == "" {
			err = fmt.Errorf("invalid proxy url")
		}
	}
	return
}
