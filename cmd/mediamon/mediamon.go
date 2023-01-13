package main

import (
	"errors"
	"fmt"
	"github.com/clambin/mediamon/collectors"
	"github.com/clambin/mediamon/services"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
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

	var opts slog.HandlerOptions
	if cfg.Debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}
	slog.SetDefault(slog.New(opts.NewTextHandler(os.Stdout)))

	if cfg.Services, err = services.ParseConfigFile(servicesFilename); err != nil {
		slog.Error("unable to parse Services file", err, "filename", servicesFilename)
		return
	}

	prometheus.MustRegister(collectors.Create(cfg.Services))
	go runPrometheusServer(cfg.Port)

	slog.Info("mediamon started", "version", version.BuildVersion)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	slog.Info("mediamon exiting")
}

func runPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus listener", err)
		panic(err)
	}
}
