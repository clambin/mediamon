package mediamon

import (
	"time"

	log "github.com/sirupsen/logrus"

	"mediamon/internal/bandwidth"
	"mediamon/internal/connectivity"
	"mediamon/internal/plex"
	"mediamon/internal/services"
	"mediamon/internal/transmission"
	"mediamon/internal/xxxarr"
)

type Configuration struct {
	Port     int
	Debug    bool
	Services *services.Config
}

func StartProbes(cfg *Configuration) []string {
	probes := make([]string, 0)

	// Transmission Probe
	if cfg.Services.Transmission.URL != "" {
		log.Debugf("Starting Transmission probe (%s)", cfg.Services.Transmission.URL)
		runProbe(
			transmission.NewProbe(cfg.Services.Transmission.URL),
			cfg.Services.Transmission.Interval,
		)
		probes = append(probes, "Transmission")
	}

	// Sonarr Probe
	if cfg.Services.Sonarr.URL != "" {
		log.Debugf("Starting Sonarr probe (%s)", cfg.Services.Sonarr.URL)
		runProbe(
			xxxarr.NewProbe(cfg.Services.Sonarr.URL, cfg.Services.Sonarr.APIKey, "sonarr"),
			cfg.Services.Sonarr.Interval,
		)
		probes = append(probes, "Sonarr")
	}

	// Radarr Probe
	if cfg.Services.Radarr.URL != "" {
		log.Debugf("Starting Radarr probe (%s)", cfg.Services.Radarr.URL)
		runProbe(
			xxxarr.NewProbe(cfg.Services.Radarr.URL, cfg.Services.Radarr.APIKey, "radarr"),
			cfg.Services.Radarr.Interval,
		)
		probes = append(probes, "Radarr")
	}

	// Plex Probe
	if cfg.Services.Plex.URL != "" {
		log.Debugf("Starting Plex probe (%s)", cfg.Services.Plex.URL)
		runProbe(
			plex.NewProbe(cfg.Services.Plex.URL, cfg.Services.Plex.UserName, cfg.Services.Plex.Password),
			cfg.Services.Plex.Interval,
		)
		probes = append(probes, "Plex")
	}

	// Bandwidth Probe
	if cfg.Services.OpenVPN.Bandwidth.FileName != "" {
		log.Debugf("Starting Bandwidth probe (%s)", cfg.Services.OpenVPN.Bandwidth.FileName)
		runProbe(
			bandwidth.NewProbe(cfg.Services.OpenVPN.Bandwidth.FileName),
			cfg.Services.OpenVPN.Bandwidth.Interval,
		)
		probes = append(probes, "Bandwidth")

	}

	// Connectivity Probe
	if cfg.Services.OpenVPN.Connectivity.ProxyURL != nil {
		log.Debugf("Starting Connectivity probe (%s)", cfg.Services.OpenVPN.Connectivity.ProxyURL)
		runProbe(
			connectivity.NewProbe(cfg.Services.OpenVPN.Connectivity.ProxyURL, cfg.Services.OpenVPN.Connectivity.Token),
			cfg.Services.OpenVPN.Connectivity.Interval,
		)
		probes = append(probes, "Connectivity")
	}

	return probes
}

// Helper to start individual probes

type runnable interface {
	Run()
}

func runProbe(probe runnable, interval time.Duration) {
	go func() {
		for {
			probe.Run()
			time.Sleep(interval)
		}
	}()
}
