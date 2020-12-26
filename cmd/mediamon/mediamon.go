package main

import (
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"mediamon/internal/bandwidth"
	"mediamon/internal/connectivity"
	"mediamon/internal/metrics"
	"mediamon/internal/plex"
	"mediamon/internal/services"
	"mediamon/internal/transmission"
	"mediamon/internal/version"
	"mediamon/internal/xxxarr"
)

func main() {
	cfg := struct {
		port             int
		debug            bool
		servicesFilename string
		services         *services.Config
	}{}

	a := kingpin.New(filepath.Base(os.Args[0]), "media monitor")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&cfg.debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&cfg.port)
	a.Flag("file", "Service configuration file").Required().StringVar(&cfg.servicesFilename)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if cfg.debug {
		log.SetLevel(log.DebugLevel)
	}

	cfg.services, err = services.ParseConfigFile(cfg.servicesFilename)

	if err != nil {
		log.Warningf("unable to parse services file '%s': %v", cfg.servicesFilename, err)
		os.Exit(3)
	}

	log.Info("media monitor v" + version.BuildVersion)

	log.Debug(cfg.services)

	// Transmission Probe
	if cfg.services.Transmission.URL != "" {
		log.Debugf("Starting Transmission probe (%s)", cfg.services.Transmission.URL)
		runProbe(
			transmission.NewProbe(cfg.services.Transmission.URL),
			cfg.services.Transmission.Interval,
		)
	}

	// Sonarr Probe
	if cfg.services.Sonarr.URL != "" {
		log.Debugf("Starting Sonarr probe (%s)", cfg.services.Sonarr.URL)
		runProbe(
			xxxarr.NewProbe(cfg.services.Sonarr.URL, cfg.services.Sonarr.APIKey, "sonarr"),
			cfg.services.Sonarr.Interval,
		)
	}

	// Radarr Probe
	if cfg.services.Radarr.URL != "" {
		log.Debugf("Starting Radarr probe (%s)", cfg.services.Radarr.URL)
		runProbe(
			xxxarr.NewProbe(cfg.services.Radarr.URL, cfg.services.Radarr.APIKey, "radarr"),
			cfg.services.Radarr.Interval,
		)
	}

	// Plex Probe
	if cfg.services.Plex.URL != "" {
		log.Debugf("Starting Plex probe (%s)", cfg.services.Plex.URL)
		runProbe(
			plex.NewProbe(cfg.services.Plex.URL, cfg.services.Plex.UserName, cfg.services.Plex.Password),
			cfg.services.Plex.Interval,
		)
	}

	// Bandwidth Probe
	if cfg.services.OpenVPN.Bandwidth.FileName != "" {
		log.Debugf("Starting Bandwidth probe (%s)", cfg.services.OpenVPN.Bandwidth.FileName)
		runProbe(
			bandwidth.NewProbe(cfg.services.OpenVPN.Bandwidth.FileName),
			cfg.services.OpenVPN.Bandwidth.Interval,
		)
	}

	// Connectivity Probe
	if cfg.services.OpenVPN.Connectivity.Proxy != "" {
		if proxyURL, err := url.Parse(cfg.services.OpenVPN.Connectivity.Proxy); err == nil {
			log.Debugf("Starting Connectivity probe (%s)", cfg.services.OpenVPN.Connectivity.Proxy)
			runProbe(
				connectivity.NewProbe(proxyURL, cfg.services.OpenVPN.Connectivity.Token),
				cfg.services.OpenVPN.Connectivity.Interval,
			)
		} else {
			log.Warningf("connectivity: invalid Proxy URL (%s - %s)",
				cfg.services.OpenVPN.Connectivity.Proxy,
				err.Error(),
			)
		}
	}

	// Prometheus Metrics
	metrics.Run(cfg.port, false)
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
