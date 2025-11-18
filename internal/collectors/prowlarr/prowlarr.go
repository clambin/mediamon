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

func newMetrics(url string) map[string]*prometheus.Desc {
	constLabels := prometheus.Labels{"application": "prowlarr", "url": url}
	return map[string]*prometheus.Desc{
		"indexerResponseTime": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_response_time"),
			"Average response time in seconds",
			[]string{"indexer"},
			constLabels,
		),
		"indexerQueryTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_query_total"),
			"Total number of queries to this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"indexerGrabTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_grab_total"),
			"Total number of grabs from this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"indexerFailedQueryTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_failed_query_total"),
			"Total number of failed queries to this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"indexerFailedGrabTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_failed_grab_total"),
			"Total number of failed grabs from this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"userAgentQueryTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "user_agent_query_total"),
			"Total number of queries by user agent",
			[]string{"user_agent"},
			constLabels,
		),
		"userAgentGrabTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "user_agent_grab_total"),
			"Total number of grabs by user agent",
			[]string{"user_agent"},
			constLabels,
		),
	}
}

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
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerResponseTime"], prometheus.GaugeValue, float64(*indexer.AverageResponseTime)/1000, name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerQueryTotal"], prometheus.CounterValue, float64(*indexer.NumberOfQueries), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedQueryTotal"], prometheus.CounterValue, float64(*indexer.NumberOfFailedQueries), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerGrabTotal"], prometheus.CounterValue, float64(*indexer.NumberOfGrabs), name)
		ch <- prometheus.MustNewConstMetric(c.metrics["indexerFailedGrabTotal"], prometheus.CounterValue, float64(*indexer.NumberOfFailedGrabs), name)
	}
	for _, userAgent := range *stats.UserAgents {
		agent := *userAgent.UserAgent
		ch <- prometheus.MustNewConstMetric(c.metrics["userAgentQueryTotal"], prometheus.CounterValue, float64(*userAgent.NumberOfQueries), agent)
		ch <- prometheus.MustNewConstMetric(c.metrics["userAgentGrabTotal"], prometheus.CounterValue, float64(*userAgent.NumberOfGrabs), agent)
	}
}
