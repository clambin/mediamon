package plex

import (
	"log/slog"
	"net/http"
	"runtime"
	"sync"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/iplocator"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	libraryCollector prometheus.Collector
	versionCollector prometheus.Collector
	logger           *slog.Logger
	sessionCollector sessionCollector
}

type Getter interface {
	identityGetter
	sessionGetter
	libraryGetter
}

type IPLocator interface {
	Locate(string) (iplocator.Location, error)
}

type Config struct {
	URL      string
	UserName string
	Password string
}

// NewCollector creates a new Collector
func NewCollector(version, url, clientID, username, password string, httpClient *http.Client, logger *slog.Logger) *Collector {
	if clientID == "" {
		clientID = uuid.New().String()
		logger.Info("clientID not set, using generated clientID", "clientID", clientID)
	}
	id := plex.ClientIdentity{
		Product:         "github.com/clambin/mediamon",
		Version:         version,
		Platform:        runtime.GOOS,
		PlatformVersion: runtime.Version(),
		DeviceName:      "Media Monitor",
		Identifier:      clientID,
	}
	p := plex.New(url,
		plex.WithHTTPClient(httpClient),
		plex.WithCredentials(username, password, id),
	)
	c := Collector{
		versionCollector: newVersionCollector(p, url, logger),
		sessionCollector: sessionCollector{
			sessionGetter: p,
			ipLocator:     iplocator.New(httpClient),
			url:           url,
			logger:        logger,
		},
		libraryCollector: newLibraryCollector(p, url, logger),
		logger:           logger,
	}
	return &c
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.versionCollector.Describe(ch)
	c.sessionCollector.Describe(ch)
	c.libraryCollector.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var g sync.WaitGroup
	g.Go(func() { c.versionCollector.Collect(ch) })
	g.Go(func() { c.sessionCollector.Collect(ch) })
	g.Go(func() { c.libraryCollector.Collect(ch) })
	g.Wait()
}
