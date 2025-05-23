package plex

import (
	"codeberg.org/clambin/go-common/httputils/metrics"
	"codeberg.org/clambin/go-common/httputils/roundtripper"
	"github.com/clambin/mediaclients/plex"
	collectorbreaker "github.com/clambin/mediamon/v2/collector-breaker"
	"github.com/clambin/mediamon/v2/iplocator"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Collector presents Plex statistics as Prometheus metrics
type Collector struct {
	versionCollector versionCollector
	sessionCollector sessionCollector
	libraryCollector libraryCollector
	metrics          metrics.RequestMetrics
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

var _ collectorbreaker.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	URL      string
	UserName string
	Password string
}

// NewCollector creates a new Collector
func NewCollector(version, url, username, password string, logger *slog.Logger) *collectorbreaker.CBCollector {
	m := metrics.NewRequestMetrics(metrics.Options{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "plex"},
		LabelValues: func(request *http.Request, i int) (method string, path string, code string) {
			r2 := chopPath(request)
			return request.Method, r2.URL.Path, strconv.Itoa(i)
		},
	})
	r := roundtripper.New(roundtripper.WithRequestMetrics(m))
	p := plex.New(username, password, "github.com/clambin/mediamon", version, url, r)
	c := Collector{
		versionCollector: versionCollector{
			identityGetter: p,
			url:            url,
			logger:         logger,
		},
		sessionCollector: sessionCollector{
			sessionGetter: p,
			ipLocator: iplocator.New(&http.Client{
				Transport: roundtripper.New(roundtripper.WithCache(roundtripper.CacheOptions{
					DefaultExpiration: 24 * time.Hour,
					CleanupInterval:   time.Hour,
				})),
			}),
			url:    url,
			logger: logger,
		},
		libraryCollector: libraryCollector{
			libraryGetter: p,
			url:           url,
			logger:        logger,
		},
		metrics: m,
		logger:  logger,
	}
	return collectorbreaker.New("plex", &c, logger)
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
