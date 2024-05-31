package collector_breaker

import (
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
	breaker   *breaker.CircuitBreaker
	cbMetrics *breaker.Metrics
	logger    *slog.Logger
}

var defaultConfiguration = breaker.Configuration{
	ErrorThreshold:   2,
	OpenDuration:     5 * time.Minute,
	SuccessThreshold: 1,
}

func New(name string, c Collector, logger *slog.Logger) *CBCollector {
	cfg := defaultConfiguration
	cfg.Metrics = breaker.NewMetrics("", "", name, nil)
	cfg.Logger = logger
	return NewWithConfiguration(c, cfg, logger)
}

func NewWithConfiguration(c Collector, cfg breaker.Configuration, logger *slog.Logger) *CBCollector {
	return &CBCollector{
		Collector: c,
		breaker:   breaker.New(cfg),
		cbMetrics: cfg.Metrics,
		logger:    logger,
	}
}

func (c *CBCollector) Describe(ch chan<- *prometheus.Desc) {
	c.Collector.Describe(ch)
	if c.cbMetrics != nil {
		c.cbMetrics.Describe(ch)
	}
}

func (c *CBCollector) Collect(ch chan<- prometheus.Metric) {
	_ = c.breaker.Do(func() (err error) {
		if err = c.Collector.CollectE(ch); err != nil {
			c.logger.Warn("collection failed", "err", err)
		}
		return err
	})

	if c.cbMetrics != nil {
		c.cbMetrics.Collect(ch)
	}
}

var _ prometheus.Collector = PassThroughCollector{}

// PassThroughCollector implements the prometheus Collector interface and passes the Collect call through
// to a collector_breaker.Collector's CollectE method.  It's used for unit testing a collector_breaker.Collector
// without getting the prometheus metrics.
type PassThroughCollector struct {
	Collector
}

func (p PassThroughCollector) Collect(metrics chan<- prometheus.Metric) {
	_ = p.CollectE(metrics)
}
