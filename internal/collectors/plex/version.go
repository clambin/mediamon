package plex

import (
	"context"
	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

var versionMetric = prometheus.NewDesc(
	prometheus.BuildFQName("mediamon", "plex", "version"),
	"version info",
	[]string{"version", "url"},
	nil,
)

var _ prometheus.Collector = versionCollector{}

type versionCollector struct {
	versionGetter
	url    string
	logger *slog.Logger
}

type versionGetter interface {
	GetIdentity(context.Context) (plex.Identity, error)
}

func (c versionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
}

func (c versionCollector) Collect(ch chan<- prometheus.Metric) {
	identity, err := c.versionGetter.GetIdentity(context.Background())
	if err != nil {
		//ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error","Error getting Plex version", nil, nil),err)
		c.logger.Error("failed to collect version", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), identity.Version, c.url)
}
