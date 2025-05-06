package transmission

import (
	"codeberg.org/clambin/go-common/httputils/metrics"
	"codeberg.org/clambin/go-common/httputils/roundtripper"
	"context"
	"fmt"
	collectorBreaker "github.com/clambin/mediamon/v2/collector-breaker"
	"github.com/hekmon/transmissionrpc/v3"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"net/url"
)

var (
	versionMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "version"),
		"version info",
		[]string{"version", "url"},
		nil,
	)

	activeTorrentsMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "active_torrent_count"),
		"Number of active torrents",
		[]string{"url"},
		nil,
	)

	pausedTorrentsMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "paused_torrent_count"),
		"Number of paused torrents",
		[]string{"url"},
		nil,
	)

	downloadSpeedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "download_speed"),
		"Transmission download speed in bytes / sec",
		[]string{"url"},
		nil,
	)

	uploadSpeedMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "transmission", "upload_speed"),
		"Transmission upload speed in bytes / sec",
		[]string{"url"},
		nil,
	)
)

// TransmissionClient interface
type TransmissionClient interface {
	SessionArgumentsGetAll(ctx context.Context) (sessionArgs transmissionrpc.SessionArguments, err error)
	SessionStats(ctx context.Context) (stats transmissionrpc.SessionStats, err error)
}

// Collector presents Transmission statistics as Prometheus metrics
type Collector struct {
	TransmissionClient TransmissionClient
	url                string
	metrics            metrics.RequestMetrics
	logger             *slog.Logger
}

var _ collectorBreaker.Collector = &Collector{}

// NewCollector creates a new Collector
func NewCollector(serverURL string, logger *slog.Logger) (*collectorBreaker.CBCollector, error) {
	c := Collector{
		url: serverURL,
		metrics: metrics.NewRequestMetrics(metrics.Options{
			Namespace:   "mediamon",
			ConstLabels: prometheus.Labels{"application": "transmission"},
		}),
		logger: logger,
	}

	ep, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid transmission server URL %q: %w", serverURL, err)
	}
	c.TransmissionClient, err = transmissionrpc.New(ep, &transmissionrpc.Config{
		CustomClient: &http.Client{Transport: roundtripper.New(roundtripper.WithRequestMetrics(c.metrics))},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating transmission client: %w", err)
	}

	return collectorBreaker.New("transmission", &c, logger), nil
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
	c.metrics.Describe(ch)
}

// CollectE implements the prometheus.Collector interface
func (c *Collector) CollectE(ch chan<- prometheus.Metric) error {
	var g errgroup.Group
	g.Go(func() error { return c.collectVersion(ch) })
	g.Go(func() error { return c.collectStats(ch) })
	err := g.Wait()
	c.metrics.Collect(ch)
	return err
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) error {
	args, err := c.TransmissionClient.SessionArgumentsGetAll(context.Background())
	if err != nil {
		return fmt.Errorf("error getting session parameters: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), *args.Version, c.url)
	return nil
}

func (c *Collector) collectStats(ch chan<- prometheus.Metric) error {
	stats, err := c.TransmissionClient.SessionStats(context.Background())
	if err != nil {
		return fmt.Errorf("error getting session statistics: %w", err)
	}
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.ActiveTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.PausedTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.DownloadSpeed), c.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.UploadSpeed), c.url)
	return nil
}
