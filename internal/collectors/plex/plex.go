package plex

import (
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/pkg/breaker"
	collectorBreaker "github.com/clambin/mediamon/v2/pkg/collector-breaker"
	"github.com/clambin/mediamon/v2/pkg/iplocator"
	customMetrics "github.com/clambin/mediamon/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	versionCollector
	sessionCollector
	libraryCollector
	metrics metrics.RequestMetrics
	logger  *slog.Logger
}

type Getter interface {
	versionGetter
	sessionGetter
	libraryGetter
}

type IPLocator interface {
	Locate(string) (float64, float64, error)
}

var breakerConfiguration = breaker.Configuration{
	FailureThreshold: 6,
	OpenDuration:     time.Minute,
	SuccessThreshold: 6,
}

var _ collectorBreaker.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	URL      string
	UserName string
	Password string
}

// NewCollector creates a new Collector
func NewCollector(version, url, username, password string, logger *slog.Logger) *collectorBreaker.CBCollector {
	m := customMetrics.NewCustomizedRoundTripMetrics("mediamon", "", map[string]string{"application": "plex"}, chopPath)
	r := roundtripper.New(roundtripper.WithRequestMetrics(m))
	p := plex.New(username, password, "github.com/clambin/mediamon", version, url, r)
	c := Collector{
		versionCollector: versionCollector{
			versionGetter: p,
			url:           url,
			logger:        logger,
		},
		sessionCollector: sessionCollector{
			sessionGetter: p,
			IPLocator:     iplocator.New(logger),
			url:           url,
			logger:        logger,
		},
		libraryCollector: libraryCollector{
			libraryGetter: p,
			url:           url,
			logger:        logger,
		},
		metrics: m,
		logger:  logger,
	}
	return collectorBreaker.New(&c, breakerConfiguration, logger)
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
	c.metrics.Describe(ch)
}

// CollectE implements the prometheus.Collector interface
func (c *Collector) CollectE(ch chan<- prometheus.Metric) error {
	var g errgroup.Group
	g.Go(func() error { return c.versionCollector.CollectE(ch) })
	g.Go(func() error { return c.sessionCollector.CollectE(ch) })
	g.Go(func() error { return c.libraryCollector.CollectE(ch) })
	err := g.Wait()
	c.metrics.Collect(ch)
	return err
}
