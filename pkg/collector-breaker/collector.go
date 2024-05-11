package collector_breaker

import (
	"github.com/clambin/mediamon/v2/pkg/breaker"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"time"
)

type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	CollectE(ch chan<- prometheus.Metric) error
}

var _ prometheus.Collector = &CBCollector{}

type CBCollector struct {
	Collector
	breaker *breaker.CircuitBreaker
}

func New(c Collector, failureThreshold int, openDuration time.Duration, successThreshold int, logger *slog.Logger) *CBCollector {
	return &CBCollector{
		Collector: c,
		breaker:   breaker.New(failureThreshold, openDuration, successThreshold, logger.With("component", "breaker")),
	}
}

func (c *CBCollector) Describe(ch chan<- *prometheus.Desc) {
	c.Collector.Describe(ch)
}

func (c *CBCollector) Collect(ch chan<- prometheus.Metric) {
	c.breaker.Do(func() error {
		return c.Collector.CollectE(ch)
	})
}
