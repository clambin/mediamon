package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"mediamon/internal/mediamon"
	"mediamon/internal/metrics"
	"os"
	"path/filepath"

	"mediamon/internal/services"
	"mediamon/internal/version"
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

	cfg.Services, err = services.ParseConfigFile(servicesFilename)

	if err != nil {
		log.Warningf("unable to parse Services file '%s': %v", servicesFilename, err)
		os.Exit(3)
	}

	log.Info("media monitor v" + version.BuildVersion)
	log.Debug(cfg.Services)

	mediamon.StartProbes(&cfg)
	metrics.Run(cfg.Port, false)
}
