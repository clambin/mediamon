package collectors

import (
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/services"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
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
		slog.Info("monitoring Transmission", "url", cfg.Transmission.URL)
		c.Transmission = transmission.NewCollector(cfg.Transmission.URL)
	}

	// Sonarr Collector
	if cfg.Sonarr.URL != "" {
		slog.Info("monitoring Sonarr", "url", cfg.Sonarr.URL)
		c.Sonarr = xxxarr.NewSonarrCollector(cfg.Sonarr.URL, cfg.Sonarr.APIKey)
	}

	// Radarr Collector
	if cfg.Radarr.URL != "" {
		slog.Info("monitoring Radarr", "url", cfg.Radarr.URL)
		c.Radarr = xxxarr.NewRadarrCollector(cfg.Radarr.URL, cfg.Radarr.APIKey)
	}

	// Plex Collector
	if cfg.Plex.URL != "" {
		slog.Info("monitoring Plex", "url", cfg.Plex.URL)
		c.Plex = plex.NewCollector(cfg.Plex.URL, cfg.Plex.UserName, cfg.Plex.Password)
	}

	// Bandwidth Probe
	if cfg.OpenVPN.Bandwidth.FileName != "" {
		slog.Info("monitoring OpenVPN Bandwidth usage", "filename", cfg.OpenVPN.Bandwidth.FileName)
		c.Bandwidth = bandwidth.NewCollector(cfg.OpenVPN.Bandwidth.FileName)
	}

	// Connectivity Probe
	if cfg.OpenVPN.Connectivity.Token != "" {
		// proxyURL has already been validated when we loaded the configuration
		proxyURL, _ := url.Parse(cfg.OpenVPN.Connectivity.Proxy)
		slog.Info("monitoring OpenVPN connectivity", "proxyURL", proxyURL.String())
		c.Connectivity = connectivity.NewCollector(cfg.OpenVPN.Connectivity.Token, proxyURL, cfg.OpenVPN.Connectivity.Interval)
	}

	return c
}
