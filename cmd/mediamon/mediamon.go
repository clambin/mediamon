package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/clambin/mediamon/internal/mediamon"
	"github.com/clambin/mediamon/internal/services"
	"github.com/clambin/mediamon/internal/version"
)

func main() {
	var (
		servicesFilename string
	)

	cfg := mediamon.Configuration{}

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

	mediamon.StartProbes(&cfg)

	// Run initialized & runs the metrics
	listenAddress := fmt.Sprintf(":%d", cfg.Port)
	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(listenAddress, nil)
	log.WithField("err", err).Error("Failed to start Prometheus http handler")
}
