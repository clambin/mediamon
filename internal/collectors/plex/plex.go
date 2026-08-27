package plex

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"sync"
	"uuid"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediaclients/plex/plextv"
	"github.com/clambin/mediaclients/plex/vault"
	"github.com/clambin/mediamon/v2/iplocator"
	"github.com/prometheus/client_golang/prometheus"
)

// Config holds the configuration for the Plex collector
type Config struct {
	Token         string
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
func NewCollector(url string, pcfg Config, httpClient *http.Client, logger *slog.Logger) (*Collector, error) {
	token := pcfg.Token
	if token == "" {
		// we don't have a token, so we need to get the PMS's token from Plex.tv

		// we need a client ID to register our client with Plex.tv. Create one if not set.
		if pcfg.ClientID == "" {
			pcfg.ClientID = uuid.New().String()
			logger.Info("clientID not set, using generated clientID", "clientID", pcfg.ClientID)
		}

		// build plex.tv Config to create a client
		config := plextv.DefaultConfig().
			WithClientID(pcfg.ClientID).
			WithDevice(plextv.DeviceInformation{
				Product:    "github.com/clambin/mediamon",
				Version:    pcfg.Version,
				DeviceName: "Media Monitor",
				Platform:   runtime.GOOS,
				Provides:   "controller",
			})
		ctx := context.Background()
		plexTVClient := config.Client(ctx, config.TokenSource(append(pcfg.options(), plextv.WithLogger(logger))...))

		// get the Plex Media Server token from plex.tv
		var err error
		if token, err = plex.Token(ctx, url, plexTVClient); err != nil {
			return nil, fmt.Errorf("failed to get Plex token: %w", err)
		}
	}

	// create a Plex client
	client := plex.New(url, token, plex.WithHTTPClient(httpClient))

	// create the collector
	c := Collector{
		collectors: []prometheus.Collector{
			newVersionCollector(client, url, logger),
			&sessionCollector{
				sessionGetter: client,
				ipLocator:     iplocator.New(httpClient),
				url:           url,
				logger:        logger,
			},
			newLibraryCollector(client, url, logger),
		},
	}
	return &c, nil
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
