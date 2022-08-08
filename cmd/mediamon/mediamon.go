package main

import (
	"github.com/clambin/go-metrics/server"
	"github.com/clambin/mediamon/collectors"
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
		log.WithError(err).WithField("filename", servicesFilename).Fatal("unable to parse Services file")
	}

	log.WithField("version", version.BuildVersion).Info("media monitor starting")

	collectors.Create(cfg.Services, prometheus.DefaultRegisterer)

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
