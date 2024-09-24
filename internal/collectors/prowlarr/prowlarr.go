package prowlarr

import (
	"context"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/mediaclients/prowlarr"
	collectorbreaker "github.com/clambin/mediamon/v2/collector-breaker"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	"time"
)

type Collector struct {
	ProwlarrClient
	metrics      map[string]*prometheus.Desc
	tpMetrics    metrics.RequestMetrics
	cacheMetrics roundtripper.CacheMetrics
	logger       *slog.Logger
}

type ProwlarrClient interface {
	GetApiV1IndexerstatsWithResponse(ctx context.Context, params *prowlarr.GetApiV1IndexerstatsParams, reqEditors ...prowlarr.RequestEditorFn) (*prowlarr.GetApiV1IndexerstatsResponse, error)
}

func New(url, apiKey string, logger *slog.Logger) *collectorbreaker.CBCollector {
	c := Collector{
		metrics: newMetrics(url),
		tpMetrics: metrics.NewRequestMetrics(metrics.Options{
			Namespace:   "mediamon",
			ConstLabels: prometheus.Labels{"application": "prowlarr"},
		}),
		cacheMetrics: roundtripper.NewCacheMetrics(roundtripper.CacheMetricsOptions{
			Namespace:   "mediamon",
			ConstLabels: prometheus.Labels{"application": "prowlarr"},
		}),
		logger: logger,
	}

	r := roundtripper.New(
		roundtripper.WithCache(roundtripper.CacheOptions{
			DefaultExpiration: 15 * time.Minute,
			CleanupInterval:   time.Hour,
			CacheMetrics:      c.cacheMetrics,
		}),
		roundtripper.WithRequestMetrics(c.tpMetrics),
	)

	var err error
	c.ProwlarrClient, err = prowlarr.NewClientWithResponses(
		url,
		prowlarr.WithRequestEditorFn(clients.WithToken(apiKey)),
		prowlarr.WithHTTPClient(&http.Client{Transport: r}),
	)
	if err != nil {
		// TODO
		panic(err)
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
	stats, err := c.ProwlarrClient.GetApiV1IndexerstatsWithResponse(context.Background(), nil)
	if err == nil {
		for _, indexer := range *stats.JSON200.Indexers {
			name := *indexer.IndexerName
			//c.logger.Debug("indexer found", "indexer", name, "queries", indexer.NumberOfQueries, "grabs", indexer.NumberOfGrabs)
			responseTimeMsec := *indexer.AverageResponseTime
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerResponseTime"], prometheus.GaugeValue, float64(responseTimeMsec)/1000, name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerQueryTotal"], prometheus.CounterValue, float64(*indexer.NumberOfQueries), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedQueryTotal"], prometheus.CounterValue, float64(*indexer.NumberOfFailedQueries), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerGrabTotal"], prometheus.CounterValue, float64(*indexer.NumberOfGrabs), name)
			ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedGrabTotal"], prometheus.CounterValue, float64(*indexer.NumberOfFailedGrabs), name)
		}
		for _, userAgent := range *stats.JSON200.UserAgents {
			agent := *userAgent.UserAgent
			//c.logger.Debug("user agent found", "agent", agent, "queries", userAgent.NumberOfQueries, "grabs", userAgent.NumberOfGrabs)
			ch <- prometheus.MustNewConstMetric(c.metrics["userAgentQueryTotal"], prometheus.CounterValue, float64(*userAgent.NumberOfQueries), agent)
			ch <- prometheus.MustNewConstMetric(c.metrics["userAgentGrabTotal"], prometheus.CounterValue, float64(*userAgent.NumberOfGrabs), agent)
		}
	}
	c.tpMetrics.Collect(ch)
	c.cacheMetrics.Collect(ch)
	return err
}
