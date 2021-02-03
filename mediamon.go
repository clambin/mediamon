package main

import (
	"fmt"
	"github.com/clambin/mediamon/internal/bandwidth"
	"github.com/clambin/mediamon/internal/connectivity"
	"github.com/clambin/mediamon/internal/plex"
	"github.com/clambin/mediamon/internal/services"
	"github.com/clambin/mediamon/internal/transmission"
	"github.com/clambin/mediamon/internal/version"
	"github.com/clambin/mediamon/internal/xxxarr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type configuration struct {
	Port     int
	Debug    bool
	Services *services.Config
}

func main() {
	var (
		servicesFilename string
	)

	cfg := configuration{}

	a := kingpin.New(filepath.Base(os.Args[0]), "media monitor")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&cfg.Debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&cfg.Port)
	a.Flag("file", "Service configuration file").Required().StringVar(&servicesFilename)

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

	log.Info("media monitor v" + version.BuildVersion)

	go func() {
		// Run initialized & runs the metrics
		listenAddress := fmt.Sprintf(":%d", cfg.Port)
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(listenAddress, nil)
		log.WithField("err", err).Error("Failed to start Prometheus http handler")
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Transmission Probe
	var transmissionProbe *transmission.Probe
	if cfg.Services.Transmission.URL != "" {
		log.WithField("url", cfg.Services.Transmission.URL).Info("monitoring Transmission")
		transmissionProbe = transmission.NewProbe(cfg.Services.Transmission.URL)
	}
	transmissionTicker := time.NewTicker(cfg.Services.Transmission.Interval)

	// Sonarr Probe
	var sonarrProbe *xxxarr.Probe
	if cfg.Services.Sonarr.URL != "" {
		log.WithField("url", cfg.Services.Sonarr.URL).Info("monitoring Sonarr")
		sonarrProbe = xxxarr.NewProbe(cfg.Services.Sonarr.URL, cfg.Services.Sonarr.APIKey, "sonarr")
	}
	sonarrTicker := time.NewTicker(cfg.Services.Sonarr.Interval)

	// Radarr Probe
	var radarrProbe *xxxarr.Probe
	if cfg.Services.Radarr.URL != "" {
		log.WithField("url", cfg.Services.Radarr.URL).Info("monitoring Radarr")
		radarrProbe = xxxarr.NewProbe(cfg.Services.Radarr.URL, cfg.Services.Radarr.APIKey, "radarr")
	}
	radarrTicker := time.NewTicker(cfg.Services.Radarr.Interval)

	// Plex Probe
	var plexProbe *plex.Probe
	if cfg.Services.Plex.URL != "" {
		log.WithField("url", cfg.Services.Plex.URL).Info("monitoring Plex")
		plexProbe = plex.NewProbe(cfg.Services.Plex.URL, cfg.Services.Plex.UserName, cfg.Services.Plex.Password)
	}
	plexTicker := time.NewTicker(cfg.Services.Plex.Interval)

	// Bandwidth Probe
	var bandwidthProbe *bandwidth.Probe
	if cfg.Services.OpenVPN.Bandwidth.FileName != "" {
		log.WithField("filename", cfg.Services.OpenVPN.Bandwidth.FileName).Info("monitoring OpenVPN Bandwidth usage")
		bandwidthProbe = bandwidth.NewProbe(cfg.Services.OpenVPN.Bandwidth.FileName)
	}
	bandwidthTicker := time.NewTicker(cfg.Services.OpenVPN.Bandwidth.Interval)

	// Connectivity Probe
	var connectivityProbe *connectivity.Probe
	if cfg.Services.OpenVPN.Connectivity.ProxyURL != nil {
		log.WithField("proxyURL", cfg.Services.OpenVPN.Connectivity.ProxyURL).Info("monitoring OpenVPN connectivity")
		connectivityProbe = connectivity.NewProbe(cfg.Services.OpenVPN.Connectivity.ProxyURL, cfg.Services.OpenVPN.Connectivity.Token)
	}
	connectivityTicker := time.NewTicker(cfg.Services.OpenVPN.Connectivity.Interval)

loop:
	for {
		select {
		case <-transmissionTicker.C:
			if transmissionProbe != nil {
				_ = transmissionProbe.Run()
			}
		case <-sonarrTicker.C:
			if sonarrProbe != nil {
				_ = sonarrProbe.Run()
			}
		case <-radarrTicker.C:
			if radarrProbe != nil {
				_ = radarrProbe.Run()
			}
		case <-plexTicker.C:
			if plexProbe != nil {
				_ = plexProbe.Run()
			}
		case <-bandwidthTicker.C:
			if bandwidthProbe != nil {
				_ = bandwidthProbe.Run()
			}
		case <-connectivityTicker.C:
			if connectivityProbe != nil {
				_ = connectivityProbe.Run()
			}
		case <-interrupt:
			break loop
		}
	}

	log.Info("mediamon exiting")
}
