package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"mediamon/internal/plex"
	"mediamon/internal/xxxarr"
	"os"
	"path/filepath"
	"time"

	"mediamon/internal/metrics"
	"mediamon/internal/services"
	"mediamon/internal/transmission"
	"mediamon/internal/version"
)

func main() {
	cfg := struct {
		port             int
		debug            bool
		interval         string
		servicesFilename string
		services         services.Config
	}{}

	a := kingpin.New(filepath.Base(os.Args[0]), "media monitor")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&cfg.debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&cfg.port)
	a.Flag("interval", "Time between measurements").Default("30s").StringVar(&cfg.interval)
	a.Flag("file", "Service configuration file").Required().StringVar(&cfg.servicesFilename)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if cfg.servicesFilename != "" {
		err = services.ParseConfigFile(cfg.servicesFilename, &cfg.services)
	}

	if cfg.debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("media monitor v" + version.BuildVersion)

	// Prometheus Metrics
	metrics.Init(cfg.port)

	// Transmission Probe
	if cfg.services.Transmission.URL != "" {

		log.Debugf("Starting Transmission probe (%s)", cfg.services.Transmission.URL)

		interval := cfg.services.Transmission.Interval
		duration, err := time.ParseDuration(interval)
		if err != nil {
			log.Warningf("Failed to parse Transmission duration '%s'. Defaulting to 30s", interval)
			duration = 30 * time.Second
		}

		go func(duration time.Duration) {
			probe := transmission.NewProbe(cfg.services.Transmission.URL)

			for {
				probe.Run()
				time.Sleep(duration)
			}
		}(duration)
	}

	// Sonarr Probe
	if cfg.services.Sonarr.URL != "" {

		log.Debugf("Starting Sonarr probe (%s)", cfg.services.Sonarr.URL)

		interval := cfg.services.Sonarr.Interval
		duration, err := time.ParseDuration(interval)
		if err != nil {
			log.Warningf("Failed to parse Sonarr duration '%s'. Defaulting to 30s", interval)
			duration = 30 * time.Second
		}

		go func(duration time.Duration) {
			probe := xxxarr.NewProbe(cfg.services.Sonarr.URL, cfg.services.Sonarr.APIKey, "sonarr")

			for {
				probe.Run()
				time.Sleep(duration)
			}
		}(duration)
	}

	// Radarr Probe
	if cfg.services.Radarr.URL != "" {

		log.Debugf("Starting Radarr probe (%s)", cfg.services.Radarr.URL)

		interval := cfg.services.Sonarr.Interval
		duration, err := time.ParseDuration(interval)
		if err != nil {
			log.Warningf("Failed to parse Radarr duration '%s'. Defaulting to 30s", interval)
			duration = 30 * time.Second
		}

		go func(duration time.Duration) {
			probe := xxxarr.NewProbe(cfg.services.Radarr.URL, cfg.services.Radarr.APIKey, "radarr")

			for {
				probe.Run()
				time.Sleep(duration)
			}
		}(duration)
	}

	// Plex Probe
	if cfg.services.Plex.URL != "" {

		log.Debugf("Starting Plex probe (%s)", cfg.services.Plex.URL)

		username := cfg.services.Plex.UserName
		password := cfg.services.Plex.Password
		interval := cfg.services.Plex.Interval
		duration, err := time.ParseDuration(interval)
		if err != nil {
			log.Warningf("Failed to parse Plex duration '%s'. Defaulting to 30s", interval)
			duration = 30 * time.Second
		}

		go func(duration time.Duration) {
			probe := plex.NewProbe(cfg.services.Plex.URL, username, password)

			for {
				probe.Run()
				time.Sleep(duration)
			}
		}(duration)
	}

	for {
		time.Sleep(30 * time.Second)
	}
}
