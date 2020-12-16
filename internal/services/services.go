package services

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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
}

func ParseConfigFile(fileName string, config *Config) error {
	content, err := ioutil.ReadFile(fileName)
	if err == nil {
		err = ParseConfig(content, config)
	}
	return err
}

func ParseConfig(content []byte, config *Config) error {
	return yaml.Unmarshal(content, config)
}
