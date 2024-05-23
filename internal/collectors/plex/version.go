package plex

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/plex"
	collector_breaker "github.com/clambin/mediamon/v2/pkg/collector-breaker"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

var versionMetric = prometheus.NewDesc(
	prometheus.BuildFQName("mediamon", "plex", "version"),
	"version info",
	[]string{"version", "url"},
	nil,
)

var _ collector_breaker.Collector = versionCollector{}

type versionCollector struct {
	identityGetter
	url    string
	logger *slog.Logger
}

type identityGetter interface {
	GetIdentity(context.Context) (plex.Identity, error)
}

func (c versionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
}

func (c versionCollector) CollectE(ch chan<- prometheus.Metric) error {
	identity, err := c.identityGetter.GetIdentity(context.Background())
	if err != nil {
		return fmt.Errorf("identity: %w", err)
	}

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.Version, c.url)
	return nil
}
