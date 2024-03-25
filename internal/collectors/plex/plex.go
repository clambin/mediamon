package plex

import (
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/pkg/iplocator"
	customMetrics "github.com/clambin/mediamon/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	URL      string
	UserName string
	Password string
}

// NewCollector creates a new Collector
func NewCollector(version, url, username, password string) *Collector {
	m := customMetrics.NewCustomizedRoundTripMetrics("mediamon", "", map[string]string{"application": "plex"}, chopPath)
	r := roundtripper.New(roundtripper.WithRequestMetrics(m))
	p := plex.New(username, password, "github.com/clambin/mediamon", version, url, r)
	l := slog.Default().With(slog.String("collector", "plex"))
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
			logger:        l,
		},
		metrics: m,
		logger:  l,
	}
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

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { defer wg.Done(); c.versionCollector.Collect(ch) }()
	go func() { defer wg.Done(); c.sessionCollector.Collect(ch) }()
	go func() { defer wg.Done(); c.libraryCollector.Collect(ch) }()
	wg.Wait()
	c.metrics.Collect(ch)
}
