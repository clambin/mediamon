package services

import (
	"errors"
	log "github.com/sirupsen/logrus"
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
