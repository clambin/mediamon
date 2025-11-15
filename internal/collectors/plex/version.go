package plex

import (
	"context"
	"log/slog"

	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus"
)

var versionMetric = prometheus.NewDesc(
	prometheus.BuildFQName("mediamon", "plex", "version"),
	"version info",
	[]string{"version", "url"},
	nil,
)

type versionCollector struct {
	identityGetter identityGetter
	logger         *slog.Logger
	url            string
}

type identityGetter interface {
	GetIdentity(context.Context) (plex.Identity, error)
}

func (c versionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
}

func (c versionCollector) Collect(ch chan<- prometheus.Metric) {
	identity, err := c.identityGetter.GetIdentity(context.Background())
	if err != nil {
		c.logger.Error("failed to get identity", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.Version, c.url)
}
