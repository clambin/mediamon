package main

import (
	"errors"
	"fmt"
	"github.com/clambin/mediamon/collectors"
	"github.com/clambin/mediamon/services"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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

	go runPrometheusServer(cfg.Port)

	log.Info("mediamon started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Info("mediamon exiting")
}

func runPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Fatal("failed to start Prometheus listener")
	}
}
