package main

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/internal/bandwidth"
	"github.com/clambin/mediamon/internal/connectivity"
	"github.com/clambin/mediamon/internal/plex"
	"github.com/clambin/mediamon/internal/scheduler"
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

	log.WithField("version", version.BuildVersion).Info("media monitor starting")

	go func() {
		// Run initialized & runs the metrics
		listenAddress := fmt.Sprintf(":%d", cfg.Port)
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(listenAddress, nil)
		log.WithField("err", err).Error("Failed to start Prometheus http handler")
	}()

	// Scheduler will run the probes at their configured interval
	ctx, cancel := context.WithCancel(context.Background())
	schedulr := scheduler.New()
	go schedulr.Run(ctx)

	// Transmission Probe
	if cfg.Services.Transmission.URL != "" {
		log.WithField("url", cfg.Services.Transmission.URL).Info("monitoring Transmission")
		schedulr.Schedule <- &scheduler.ScheduledTask{
			Task:     transmission.NewProbe(cfg.Services.Transmission.URL),
			Interval: cfg.Services.Transmission.Interval,
		}
	}

	// Sonarr Probe
	if cfg.Services.Sonarr.URL != "" {
		log.WithField("url", cfg.Services.Sonarr.URL).Info("monitoring Sonarr")
		schedulr.Schedule <- &scheduler.ScheduledTask{
			Task:     xxxarr.NewProbe(cfg.Services.Sonarr.URL, cfg.Services.Sonarr.APIKey, "sonarr"),
			Interval: cfg.Services.Sonarr.Interval,
		}
	}

	// Radarr Probe
	if cfg.Services.Radarr.URL != "" {
		log.WithField("url", cfg.Services.Radarr.URL).Info("monitoring Radarr")
		schedulr.Schedule <- &scheduler.ScheduledTask{
			Task:     xxxarr.NewProbe(cfg.Services.Radarr.URL, cfg.Services.Radarr.APIKey, "radarr"),
			Interval: cfg.Services.Radarr.Interval,
		}
	}

	// Plex Probe
	if cfg.Services.Plex.URL != "" {
		log.WithField("url", cfg.Services.Plex.URL).Info("monitoring Plex")
		schedulr.Schedule <- &scheduler.ScheduledTask{
			Task:     plex.NewProbe(cfg.Services.Plex.URL, cfg.Services.Plex.UserName, cfg.Services.Plex.Password),
			Interval: cfg.Services.Plex.Interval,
		}
	}

	// Bandwidth Probe
	if cfg.Services.OpenVPN.Bandwidth.FileName != "" {
		log.WithField("filename", cfg.Services.OpenVPN.Bandwidth.FileName).Info("monitoring OpenVPN Bandwidth usage")
		schedulr.Schedule <- &scheduler.ScheduledTask{
			Task:     bandwidth.NewProbe(cfg.Services.OpenVPN.Bandwidth.FileName),
			Interval: cfg.Services.OpenVPN.Bandwidth.Interval,
		}
	}

	// Connectivity Probe
	if cfg.Services.OpenVPN.Connectivity.ProxyURL != nil {
		log.WithField("proxyURL", cfg.Services.OpenVPN.Connectivity.ProxyURL).Info("monitoring OpenVPN connectivity")
		schedulr.Schedule <- &scheduler.ScheduledTask{
			Task:     connectivity.NewProbe(cfg.Services.OpenVPN.Connectivity.ProxyURL, cfg.Services.OpenVPN.Connectivity.Token),
			Interval: cfg.Services.OpenVPN.Connectivity.Interval,
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	cancel()
	log.Info("mediamon exiting")
}
