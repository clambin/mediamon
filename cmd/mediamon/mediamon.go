package main

import (
	"github.com/clambin/go-metrics/server"
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/services"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type configuration struct {
	Port     int
	Debug    bool
	Services *services.Config
}

func main() {
	var servicesFilename string

	cfg := configuration{}

	a := kingpin.New(filepath.Base(os.Args[0]), "media monitor")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&cfg.Debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&cfg.Port)
	a.Flag("file", "Service configuration file").Required().ExistingFileVar(&servicesFilename)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if cfg.Services, err = services.ParseConfigFile(servicesFilename); err != nil {
		log.WithFields(log.Fields{"err": err, "filename": servicesFilename}).Warning("unable to parse Services file")
		os.Exit(3)
	}

	log.WithField("version", version.BuildVersion).Info("media monitor starting")

	startCollectors(&cfg)

	s := server.New(cfg.Port)
	go func() {
		err = s.Run()
		if err != http.ErrServerClosed {
			log.WithField("err", err).Error("Failed to start Prometheus http handler")
		}
	}()
	log.Info("mediamon started")

	<-shutdown.Chan()

	_ = s.Shutdown(30 * time.Second)

	log.Info("mediamon exiting")
}

func startCollectors(cfg *configuration) {
	// Transmission Collector
	if cfg.Services.Transmission.URL != "" {
		log.WithField("url", cfg.Services.Transmission.URL).Info("monitoring Transmission")
		prometheus.DefaultRegisterer.MustRegister(transmission.NewCollector(
			cfg.Services.Transmission.URL,
		))
	}

	// Sonarr Collector
	if cfg.Services.Sonarr.URL != "" {
		log.WithField("url", cfg.Services.Sonarr.URL).Info("monitoring Sonarr")
		prometheus.DefaultRegisterer.MustRegister(xxxarr.NewSonarrCollector(
			cfg.Services.Sonarr.URL,
			cfg.Services.Sonarr.APIKey,
		))
	}

	// Radarr Collector
	if cfg.Services.Radarr.URL != "" {
		log.WithField("url", cfg.Services.Radarr.URL).Info("monitoring Radarr")
		prometheus.DefaultRegisterer.MustRegister(xxxarr.NewRadarrCollector(
			cfg.Services.Radarr.URL,
			cfg.Services.Radarr.APIKey,
		))
	}

	// Plex Collector
	if cfg.Services.Plex.URL != "" {
		log.WithField("url", cfg.Services.Plex.URL).Info("monitoring Plex")
		prometheus.DefaultRegisterer.MustRegister(plex.NewCollector(
			cfg.Services.Plex.URL,
			cfg.Services.Plex.UserName,
			cfg.Services.Plex.Password,
		))
	}

	// Bandwidth Probe
	if cfg.Services.OpenVPN.Bandwidth.FileName != "" {
		log.WithField("filename", cfg.Services.OpenVPN.Bandwidth.FileName).Info("monitoring OpenVPN Bandwidth usage")
		prometheus.DefaultRegisterer.MustRegister(bandwidth.NewCollector(
			cfg.Services.OpenVPN.Bandwidth.FileName,
		))
	}

	// Connectivity Probe
	if cfg.Services.OpenVPN.Connectivity.Token != "" {
		log.WithField("proxyURL", cfg.Services.OpenVPN.Connectivity.ProxyURL).Info("monitoring OpenVPN connectivity")
		prometheus.DefaultRegisterer.MustRegister(connectivity.NewCollector(
			cfg.Services.OpenVPN.Connectivity.Token,
			cfg.Services.OpenVPN.Connectivity.ProxyURL,
			cfg.Services.OpenVPN.Connectivity.Interval,
		))
	}
}
