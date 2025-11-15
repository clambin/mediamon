package plex

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/iplocator"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	versionCollector versionCollector
	sessionCollector sessionCollector
	libraryCollector prometheus.Collector
	logger           *slog.Logger
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
func NewCollector(version, url, username, password string, httpClient *http.Client, logger *slog.Logger) *Collector {
	p := plex.New(username, password, "github.com/clambin/mediamon", version, url, httpClient.Transport)
	c := Collector{
		versionCollector: versionCollector{
			identityGetter: p,
			url:            url,
			logger:         logger,
		},
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

func chopPath(r *http.Request) *http.Request {
	path := r.URL.Path
	for _, prefix := range []string{"/library/metadata", "/library/sections"} {
		if strings.HasPrefix(path, prefix) {
			path = prefix
			break
		}
	}

	return &http.Request{
		Method: r.Method,
		URL:    &url.URL{Path: path},
	}
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
