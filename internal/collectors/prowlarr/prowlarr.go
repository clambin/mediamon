package prowlarr

import (
	"context"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/xxxarr"
	collectorbreaker "github.com/clambin/mediamon/v2/pkg/collector-breaker"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"time"
)

type Collector struct {
	client       Client
	metrics      map[string]*prometheus.Desc
	tpMetrics    metrics.RequestMetrics
	cacheMetrics roundtripper.CacheMetrics
	logger       *slog.Logger
}

type Client interface {
	GetIndexStats(context.Context) (xxxarr.ProwlarrIndexersStats, error)
}

func New(url, apiKey string, logger *slog.Logger) *collectorbreaker.CBCollector {
	tpMetrics := metrics.NewRequestMetrics(metrics.Options{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "prowlarr"},
	})
	cacheMetrics := roundtripper.NewCacheMetrics(roundtripper.CacheMetricsOptions{
		Namespace:   "mediamon",
		ConstLabels: prometheus.Labels{"application": "prowlarr"},
	})

	r := roundtripper.New(
		roundtripper.WithCache(roundtripper.CacheOptions{
			DefaultExpiration: 15 * time.Minute,
			CleanupInterval:   time.Hour,
			CacheMetrics:      cacheMetrics,
		}),
		roundtripper.WithRequestMetrics(tpMetrics),
	)
	c := Collector{
		client:       xxxarr.NewProwlarrClient(url, apiKey, r),
		metrics:      newMetrics(url),
		tpMetrics:    tpMetrics,
		cacheMetrics: cacheMetrics,
		logger:       logger,
	}
	return collectorbreaker.New("prowlarr", &c, logger)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
	c.tpMetrics.Describe(ch)
	c.cacheMetrics.Describe(ch)
}

func (c *Collector) CollectE(ch chan<- prometheus.Metric) error {
	stats, err := c.client.GetIndexStats(context.Background())
	if err == nil {
		for _, indexer := range stats.Indexers {
			name := indexer.IndexerName
			c.logger.Debug("indexer found", "indexer", name, "queries", indexer.NumberOfQueries, "grabs", indexer.NumberOfGrabs)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerResponseTime"], prometheus.GaugeValue, time.Duration(indexer.AverageResponseTime).Seconds(), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerQueryTotal"], prometheus.CounterValue, float64(indexer.NumberOfQueries), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedQueryTotal"], prometheus.CounterValue, float64(indexer.NumberOfFailedQueries), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerGrabTotal"], prometheus.CounterValue, float64(indexer.NumberOfGrabs), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedGrabTotal"], prometheus.CounterValue, float64(indexer.NumberOfFailedGrabs), name)
		}
		for _, userAgent := range stats.UserAgents {
			agent := userAgent.UserAgent
			c.logger.Debug("user agent found", "agent", agent, "queries", userAgent.NumberOfQueries, "grabs", userAgent.NumberOfGrabs)
			ch <- prometheus.MustNewConstMetric(c.metrics["userAgentQueryTotal"], prometheus.CounterValue, float64(userAgent.NumberOfQueries), agent)
			ch <- prometheus.MustNewConstMetric(c.metrics["userAgentGrabTotal"], prometheus.CounterValue, float64(userAgent.NumberOfGrabs), agent)
		}
	}
	c.tpMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
	return err
}
