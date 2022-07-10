package services

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"time"
)

// Config contains the different possible services for github.com/clambin/mediamon to monitor
type Config struct {
	Transmission struct {
		URL string
	}
	Sonarr struct {
		URL    string
		APIKey string
	}
	Radarr struct {
		URL    string
		APIKey string
	}
	Plex struct {
		URL      string
		UserName string
		Password string
	}
	OpenVPN struct {
		Bandwidth struct {
			FileName string
		}
		Connectivity struct {
			Proxy    string
			Token    string
			Interval time.Duration
		}
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
