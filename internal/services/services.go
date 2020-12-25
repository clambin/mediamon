package services

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

// Config contains the different possible services for mediamon to monitor
type Config struct {
	Transmission struct {
		URL      string
		Interval string
	}
	Sonarr struct {
		URL      string
		APIKey   string
		Interval string
	}
	Radarr struct {
		URL      string
		APIKey   string
		Interval string
	}
	Plex struct {
		URL      string
		UserName string
		Password string
		Interval string
	}
	OpenVPN struct {
		Bandwidth struct {
			FileName string
			Interval time.Duration
		}
		Connectivity struct {
			Proxy    string
			Token    string
			Interval time.Duration
		}
	}
}

// ParseConfigFile reads the configuration from the specified yaml file
func ParseConfigFile(fileName string, config *Config) error {
	content, err := ioutil.ReadFile(fileName)
	if err == nil {
		err = ParseConfig(content, config)
	}
	return err
}

// ParseConfig reads the configuration from an in-memory buffer
func ParseConfig(content []byte, config *Config) error {
	return yaml.Unmarshal(content, config)
}
