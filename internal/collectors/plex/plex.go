package plex

import (
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/roundtripper"
	"github.com/clambin/mediamon/v2/pkg/iplocator"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"time"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	versionCollector
	sessionCollector
	libraryCollector
	transport *httpclient.RoundTripper
	logger    *slog.Logger
}

type Getter interface {
	versionGetter
	sessionGetter
	libraryGetter
}

type IPLocator interface {
	Locate(string) (float64, float64, error)
}

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	URL      string
	UserName string
	Password string
}

var plexCacheTable = httpclient.CacheTable{
	{
		Path:     "/library/metadata/[0-9]+/children",
		IsRegExp: true,
		Expiry:   time.Hour,
	},
}

// NewCollector creates a new Collector
func NewCollector(version, url, username, password string) *Collector {
	r := httpclient.NewRoundTripper(
		httpclient.WithCache(plexCacheTable, time.Hour, 2*time.Hour),
		httpclient.WithCustomMetrics(roundtripper.NewRequestMeasurer("mediamon", "", "plex")),
	)
	p := plex.New(username, password, "github.com/clambin/mediamon", version, url, r)
	l := slog.Default().With("collector", "plex")
	return &Collector{
		versionCollector: versionCollector{
			versionGetter: p,
			url:           url,
			logger:        l,
		},
		sessionCollector: sessionCollector{
			sessionGetter: p,
			IPLocator:     iplocator.New(l),
			url:           url,
			logger:        l,
		},
		libraryCollector: libraryCollector{
			libraryGetter: p,
			url:           url,
			l:             l,
		},
		transport: r,
		logger:    l,
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.versionCollector.Describe(ch)
	c.sessionCollector.Describe(ch)
	c.libraryCollector.Describe(ch)
	c.transport.Describe(ch)
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	c.versionCollector.Collect(ch)
	c.sessionCollector.Collect(ch)
	c.libraryCollector.Collect(ch)
	c.transport.Collect(ch)
	c.logger.Debug("stats collected", "duration", time.Since(start))
}
