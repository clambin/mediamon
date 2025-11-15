package prowlarr

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/clambin/mediaclients/prowlarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr"
	"github.com/clambin/mediamon/v2/internal/measurer"
	"github.com/prometheus/client_golang/prometheus"
)

const refreshInterval = 15 * time.Minute

type Collector struct {
	metrics      map[string]*prometheus.Desc
	logger       *slog.Logger
	indexerStats measurer.Cached[*prowlarr.IndexerStatsResource]
}

type ProwlarrClient interface {
	GetApiV1IndexerstatsWithResponse(ctx context.Context, params *prowlarr.GetApiV1IndexerstatsParams, reqEditors ...prowlarr.RequestEditorFn) (*prowlarr.GetApiV1IndexerstatsResponse, error)
}

func New(url, apiKey string, httpClient *http.Client, logger *slog.Logger) (prometheus.Collector, error) {
	prowlarrClient, err := prowlarr.NewClientWithResponses(
		url,
		prowlarr.WithRequestEditorFn(xxxarr.WithToken(apiKey)),
		prowlarr.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating prowlarr client: %w", err)
	}

	return &Collector{
		indexerStats: measurer.Cached[*prowlarr.IndexerStatsResource]{
			Interval: refreshInterval,
			Do: func(ctx context.Context) (*prowlarr.IndexerStatsResource, error) {
				resp, err := prowlarrClient.GetApiV1IndexerstatsWithResponse(context.Background(), nil)
				if err != nil {
					return nil, fmt.Errorf("prowlarr: %w", err)
				}
				return resp.JSON200, nil
			},
		},
		metrics: newMetrics(url),
		logger:  logger,
	}, nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.indexerStats.Measure(context.Background())
	if err != nil {
		c.logger.Error("failed to get indexer stats", "err", err)
		return
	}
	for _, indexer := range *stats.Indexers {
		name := *indexer.IndexerName
		//c.logger.Debug("indexer found", "indexer", name, "queries", indexer.NumberOfQueries, "grabs", indexer.NumberOfGrabs)
		responseTimeMsec := *indexer.AverageResponseTime
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerResponseTime"], prometheus.GaugeValue, float64(responseTimeMsec)/1000, name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerQueryTotal"], prometheus.CounterValue, float64(*indexer.NumberOfQueries), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedQueryTotal"], prometheus.CounterValue, float64(*indexer.NumberOfFailedQueries), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerGrabTotal"], prometheus.CounterValue, float64(*indexer.NumberOfGrabs), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedGrabTotal"], prometheus.CounterValue, float64(*indexer.NumberOfFailedGrabs), name)
	}
	for _, userAgent := range *stats.UserAgents {
		agent := *userAgent.UserAgent
		//c.logger.Debug("user agent found", "agent", agent, "queries", userAgent.NumberOfQueries, "grabs", userAgent.NumberOfGrabs)
		ch <- prometheus.MustNewConstMetric(c.metrics["userAgentQueryTotal"], prometheus.CounterValue, float64(*userAgent.NumberOfQueries), agent)
		ch <- prometheus.MustNewConstMetric(c.metrics["userAgentGrabTotal"], prometheus.CounterValue, float64(*userAgent.NumberOfGrabs), agent)
	}
}
