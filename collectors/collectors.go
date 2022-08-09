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

// Create builds the list of collectors based on the provided configuration and registers them with a registerer.
func Create(cfg *services.Config, registerer prometheus.Registerer) (c Collectors) {
	var toRegister []prometheus.Collector

	// Transmission Collector
	if cfg.Transmission.URL != "" {
		log.WithField("url", cfg.Transmission.URL).Info("monitoring Transmission")
		c.Transmission = transmission.NewCollector(cfg.Transmission.URL)
		toRegister = append(toRegister, c.Transmission)
	}

	// Sonarr Collector
	if cfg.Sonarr.URL != "" {
		log.WithField("url", cfg.Sonarr.URL).Info("monitoring Sonarr")
		c.Sonarr = xxxarr.NewSonarrCollector(
			cfg.Sonarr.URL,
			cfg.Sonarr.APIKey,
		)
		toRegister = append(toRegister, c.Sonarr)
	}

	// Radarr Collector
	if cfg.Radarr.URL != "" {
		log.WithField("url", cfg.Radarr.URL).Info("monitoring Radarr")
		c.Radarr = xxxarr.NewRadarrCollector(
			cfg.Radarr.URL,
			cfg.Radarr.APIKey,
		)
		toRegister = append(toRegister, c.Radarr)
	}

	// Plex Collector
	if cfg.Plex.URL != "" {
		log.WithField("url", cfg.Plex.URL).Info("monitoring Plex")
		c.Plex = plex.NewCollector(
			cfg.Plex.URL,
			cfg.Plex.UserName,
			cfg.Plex.Password,
		)
		toRegister = append(toRegister, c.Plex)
	}

	// Bandwidth Probe
	if cfg.OpenVPN.Bandwidth.FileName != "" {
		log.WithField("filename", cfg.OpenVPN.Bandwidth.FileName).Info("monitoring OpenVPN Bandwidth usage")
		c.Bandwidth = bandwidth.NewCollector(
			cfg.OpenVPN.Bandwidth.FileName,
		)
		toRegister = append(toRegister, c.Bandwidth)
	}

	// Connectivity Probe
	if cfg.OpenVPN.Connectivity.Token != "" {
		// proxyURL has already been validated when we loaded the configuration
		proxyURL, _ := url.Parse(cfg.OpenVPN.Connectivity.Proxy)
		log.WithField("proxyURL", proxyURL).Info("monitoring OpenVPN connectivity")
		c.Connectivity = connectivity.NewCollector(
			cfg.OpenVPN.Connectivity.Token,
			proxyURL,
			cfg.OpenVPN.Connectivity.Interval,
		)
		toRegister = append(toRegister, c.Connectivity)
	}

	// register all collectors
	registerer.MustRegister(toRegister...)

	return
}
