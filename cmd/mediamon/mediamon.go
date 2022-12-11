package main

import (
	"errors"
	"github.com/clambin/go-common/httpserver"
	"github.com/clambin/mediamon/collectors"
	"github.com/clambin/mediamon/services"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
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
	a.Flag("port", "Prometheus metrics port").Default("9090").IntVar(&cfg.Port)
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
		log.WithError(err).WithField("filename", servicesFilename).Fatal("unable to parse Services file")
	}

	log.WithField("version", version.BuildVersion).Info("media monitor starting")

	prometheus.MustRegister(collectors.Create(cfg.Services))

	s, err := httpserver.New(
		httpserver.WithPort{Port: cfg.Port},
		httpserver.WithPrometheus{Path: "/metrics"},
	)

	if err != nil {
		log.WithError(err).Fatal("failed to start metrics server")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = s.Serve(); !errors.Is(err, http.ErrServerClosed) {
			log.WithField("err", err).Fatal("Failed to start Prometheus http handler")
		}
	}()
	log.Info("mediamon started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	_ = s.Shutdown(30 * time.Second)
	wg.Wait()
	log.Info("mediamon exiting")
}
