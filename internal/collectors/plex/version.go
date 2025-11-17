package plex

import (
	"context"
	"log/slog"
	"time"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/measurer"
	"github.com/prometheus/client_golang/prometheus"
)

const identityRefreshInterval = 15 * time.Minute

var versionMetric = prometheus.NewDesc(
	prometheus.BuildFQName("mediamon", "plex", "version"),
	"version info",
	[]string{"version", "url"},
	nil,
)

type identityGetter interface {
	GetIdentity(context.Context) (plex.Identity, error)
}

type versionCollector struct {
	identityGetter
	logger *slog.Logger
	url    string
	measurer.Cached[plex.Identity]
}

func newVersionCollector(client identityGetter, url string, logger *slog.Logger) prometheus.Collector {
	c := versionCollector{
		identityGetter: client,
		url:            url,
		logger:         logger,
	}
	c.Cached = measurer.Cached[plex.Identity]{
		Interval: identityRefreshInterval,
		Do:       func(ctx context.Context) (plex.Identity, error) { return c.GetIdentity(ctx) },
	}
	return &c
}

func (c *versionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
}

func (c *versionCollector) Collect(ch chan<- prometheus.Metric) {
	identity, err := c.Measure(context.Background())
	if err != nil {
		c.logger.Error("failed to get identity", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.Version, c.url)
}
