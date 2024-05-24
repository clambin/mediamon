package collector_breaker

import (
	"errors"
	"github.com/clambin/breaker"
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
	logger  *slog.Logger
}

var defaultConfiguration = breaker.Configuration{
	FailureThreshold: 2,
	OpenDuration:     5 * time.Minute,
	SuccessThreshold: 1,
}

func New(c Collector, logger *slog.Logger) *CBCollector {
	return NewWithConfiguration(c, defaultConfiguration, logger)
}

func NewWithConfiguration(c Collector, cfg breaker.Configuration, logger *slog.Logger) *CBCollector {
	return &CBCollector{
		Collector: c,
		breaker:   breaker.New(cfg),
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
