package transmission

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/prometheus/client_golang/prometheus"
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

type TransmissionClient interface {
	SessionArgumentsGetAll(ctx context.Context) (sessionArgs transmissionrpc.SessionArguments, err error)
	SessionStats(ctx context.Context) (stats transmissionrpc.SessionStats, err error)
}

type Collector struct {
	transmissionClient TransmissionClient
	logger             *slog.Logger
	url                string
}

// NewCollector creates a new Collector
func NewCollector(httpClient *http.Client, serverURL string, logger *slog.Logger) (prometheus.Collector, error) {
	c := Collector{
		url:    serverURL,
		logger: logger,
	}

	ep, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid transmission server URL %q: %w", serverURL, err)
	}
	c.transmissionClient, err = transmissionrpc.New(ep, &transmissionrpc.Config{CustomClient: httpClient})
	if err != nil {
		return nil, fmt.Errorf("error creating transmission client: %w", err)
	}

	return &c, nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- versionMetric
	ch <- activeTorrentsMetric
	ch <- pausedTorrentsMetric
	ch <- downloadSpeedMetric
	ch <- uploadSpeedMetric
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var g sync.WaitGroup
	g.Go(func() { c.collectVersion(ch) })
	g.Go(func() { c.collectStats(ch) })
	g.Wait()
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) {
	args, err := c.transmissionClient.SessionArgumentsGetAll(context.Background())
	if err != nil {
		c.logger.Error("error getting session parameters", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(versionMetric, prometheus.GaugeValue, float64(1), *args.Version, c.url)
}

func (c *Collector) collectStats(ch chan<- prometheus.Metric) {
	stats, err := c.transmissionClient.SessionStats(context.Background())
	if err != nil {
		c.logger.Error("error getting session statistics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(activeTorrentsMetric, prometheus.GaugeValue, float64(stats.ActiveTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(pausedTorrentsMetric, prometheus.GaugeValue, float64(stats.PausedTorrentCount), c.url)
	ch <- prometheus.MustNewConstMetric(downloadSpeedMetric, prometheus.GaugeValue, float64(stats.DownloadSpeed), c.url)
	ch <- prometheus.MustNewConstMetric(uploadSpeedMetric, prometheus.GaugeValue, float64(stats.UploadSpeed), c.url)
}
