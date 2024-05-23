package collector_breaker

import (
	"errors"
	"github.com/clambin/mediamon/v2/pkg/breaker"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	CollectE(ch chan<- prometheus.Metric) error
}

var _ prometheus.Collector = &CBCollector{}

type CBCollector struct {
	Collector
	breaker *breaker.CircuitBreaker
	logger  *slog.Logger
}

func New(c Collector, cfg breaker.Configuration, logger *slog.Logger) *CBCollector {
	return &CBCollector{
		Collector: c,
		breaker:   &breaker.CircuitBreaker{Configuration: cfg},
		logger:    logger,
	}
}

func (c *CBCollector) Describe(ch chan<- *prometheus.Desc) {
	c.Collector.Describe(ch)
}

func (c *CBCollector) Collect(ch chan<- prometheus.Metric) {
	err := c.breaker.Do(func() error {
		return c.Collector.CollectE(ch)
	})

	if err != nil && !errors.Is(err, breaker.ErrCircuitOpen) {
		c.logger.Warn("collection failed", "err", err)
	}
}
