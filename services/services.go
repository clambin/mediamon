package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"time"
)

// Config contains the different possible services for github.com/clambin/mediamon to monitor
type Config struct {
	Transmission struct {
		URL      string
		Interval time.Duration
	}
	Sonarr struct {
		URL      string
		APIKey   string
		Interval time.Duration
	}
	Radarr struct {
		URL      string
		APIKey   string
		Interval time.Duration
	}
	Plex struct {
		URL      string
		UserName string
		Password string
		Interval time.Duration
	}
	OpenVPN struct {
		Bandwidth struct {
			FileName string
			Interval time.Duration
		}
		Connectivity struct {
			Proxy    string
			ProxyURL *url.URL
			Token    string
			Interval time.Duration
		}
	}
}

// ParseConfigFile reads the configuration from the specified yaml file
func ParseConfigFile(fileName string) (*Config, error) {
	var (
		err     error
		content []byte
		config  *Config
	)
	if content, err = os.ReadFile(fileName); err == nil {
		config, err = ParseConfig(content)
	}
	return config, err
}

// ParseConfig reads the configuration from an in-memory buffer
func ParseConfig(content []byte) (*Config, error) {
	var err error

	config := Config{}
	config.Transmission.Interval = 30 * time.Second
	config.Sonarr.Interval = 30 * time.Second
	config.Radarr.Interval = 30 * time.Second
	config.Plex.Interval = 30 * time.Second
	config.OpenVPN.Bandwidth.Interval = 30 * time.Second
	config.OpenVPN.Connectivity.Interval = 5 * time.Minute

	if err = yaml.Unmarshal(content, &config); err == nil {
		if config.OpenVPN.Connectivity.Proxy != "" {
			if config.OpenVPN.Connectivity.ProxyURL, err = url.Parse(config.OpenVPN.Connectivity.Proxy); err == nil {
				if config.OpenVPN.Connectivity.ProxyURL.Scheme == "" ||
					config.OpenVPN.Connectivity.ProxyURL.Host == "" {
					err = errors.New("proxy URL is invalid")
				}
			}
		}
	}
	log.WithField("err", err).Debug("ParseConfig")

	return &config, err
}
