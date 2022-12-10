package collectors

import (
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/services"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/url"
)

// Collectors groups the different collectors
type Collectors struct {
	Transmission *transmission.Collector
	Sonarr       *xxxarr.Collector
	Radarr       *xxxarr.Collector
	Plex         *plex.Collector
	Bandwidth    *bandwidth.Collector
	Connectivity *connectivity.Collector
}

var _ prometheus.Collector = &Collectors{}

func (c Collectors) Describe(ch chan<- *prometheus.Desc) {
	if c.Transmission != nil {
		c.Transmission.Describe(ch)
	}
	if c.Sonarr != nil {
		c.Sonarr.Describe(ch)
	}
	if c.Radarr != nil {
		c.Radarr.Describe(ch)
	}
	if c.Plex != nil {
		c.Plex.Describe(ch)
	}
	if c.Bandwidth != nil {
		c.Bandwidth.Describe(ch)
	}
	if c.Connectivity != nil {
		c.Connectivity.Describe(ch)
	}
}

func (c Collectors) Collect(ch chan<- prometheus.Metric) {
	if c.Transmission != nil {
		c.Transmission.Collect(ch)
	}
	if c.Sonarr != nil {
		c.Sonarr.Collect(ch)
	}
	if c.Radarr != nil {
		c.Radarr.Collect(ch)
	}
	if c.Plex != nil {
		c.Plex.Collect(ch)
	}
	if c.Bandwidth != nil {
		c.Bandwidth.Collect(ch)
	}
	if c.Connectivity != nil {
		c.Connectivity.Collect(ch)
	}
}

// Create builds the list of collectors based on the provided configuration and registers them with a registerer.
func Create(cfg *services.Config) Collectors {
	var c Collectors

	// Transmission Collector
	if cfg.Transmission.URL != "" {
		log.WithField("url", cfg.Transmission.URL).Info("monitoring Transmission")
		c.Transmission = transmission.NewCollector(cfg.Transmission.URL)
	}

	// Sonarr Collector
	if cfg.Sonarr.URL != "" {
		log.WithField("url", cfg.Sonarr.URL).Info("monitoring Sonarr")
		c.Sonarr = xxxarr.NewSonarrCollector(cfg.Sonarr.URL, cfg.Sonarr.APIKey)
	}

	// Radarr Collector
	if cfg.Radarr.URL != "" {
		log.WithField("url", cfg.Radarr.URL).Info("monitoring Radarr")
		c.Radarr = xxxarr.NewRadarrCollector(cfg.Radarr.URL, cfg.Radarr.APIKey)
	}

	// Plex Collector
	if cfg.Plex.URL != "" {
		log.WithField("url", cfg.Plex.URL).Info("monitoring Plex")
		c.Plex = plex.NewCollector(cfg.Plex.URL, cfg.Plex.UserName, cfg.Plex.Password)
	}

	// Bandwidth Probe
	if cfg.OpenVPN.Bandwidth.FileName != "" {
		log.WithField("filename", cfg.OpenVPN.Bandwidth.FileName).Info("monitoring OpenVPN Bandwidth usage")
		c.Bandwidth = bandwidth.NewCollector(cfg.OpenVPN.Bandwidth.FileName)
	}

	// Connectivity Probe
	if cfg.OpenVPN.Connectivity.Token != "" {
		// proxyURL has already been validated when we loaded the configuration
		proxyURL, _ := url.Parse(cfg.OpenVPN.Connectivity.Proxy)
		log.WithField("proxyURL", proxyURL).Info("monitoring OpenVPN connectivity")
		c.Connectivity = connectivity.NewCollector(cfg.OpenVPN.Connectivity.Token, proxyURL, cfg.OpenVPN.Connectivity.Interval)
	}

	return c
}
