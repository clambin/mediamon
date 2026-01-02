package plex

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"sync"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediaclients/plex/plextv"
	"github.com/clambin/mediaclients/plex/vault"
	"github.com/clambin/mediamon/v2/iplocator"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

// Config holds the configuration for the Plex collector
type Config struct {
	UserName      string
	Password      string
	ClientID      string
	JWTLocation   string
	JWTPassphrase string
	Version       string
	UseJWT        bool
}

func (p Config) options() []plextv.TokenSourceOption {
	var opts []plextv.TokenSourceOption
	opts = append(opts, plextv.WithCredentials(p.UserName, p.Password))
	if p.UseJWT {
		opts = append(opts, plextv.WithJWT(vault.New[plextv.JWTSecureData](p.JWTLocation, p.JWTPassphrase)))
	}
	return opts
}

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	collectors []prometheus.Collector
}

type Getter interface {
	identityGetter
	sessionGetter
	libraryGetter
}

type IPLocator interface {
	Locate(string) (iplocator.Location, error)
}

// NewCollector creates a new Collector
func NewCollector(url string, pcfg Config, httpClient *http.Client, logger *slog.Logger) *Collector {
	if pcfg.ClientID == "" {
		pcfg.ClientID = uuid.New().String()
		logger.Info("clientID not set, using generated clientID", "clientID", pcfg.ClientID)
	}

	config := plextv.DefaultConfig().
		WithClientID(pcfg.ClientID).
		WithDevice(plextv.Device{
			Product:         "github.com/clambin/mediamon",
			Version:         pcfg.Version,
			Platform:        runtime.GOOS,
			PlatformVersion: runtime.Version(),
			DeviceName:      "Media Monitor",
			Device:          "Media Monitor Y",
			Provides:        "controller",
		})
	plexTVClient := config.Client(context.Background(), config.TokenSource(append(pcfg.options(), plextv.WithLogger(logger))...))
	pmsClient := plex.NewPMSClient(url, plexTVClient, plex.WithHTTPClient(httpClient))
	c := Collector{
		collectors: []prometheus.Collector{
			newVersionCollector(pmsClient, url, logger),
			&sessionCollector{
				sessionGetter: pmsClient,
				ipLocator:     iplocator.New(httpClient),
				url:           url,
				logger:        logger,
			},
			newLibraryCollector(pmsClient, url, logger),
		},
	}
	return &c
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var g sync.WaitGroup
	for _, collector := range c.collectors {
		g.Go(func() { collector.Collect(ch) })
	}
	g.Wait()
}
