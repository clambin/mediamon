package connectivity

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/clambin/mediamon/v2/internal/measurer"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	upMetric = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "status"),
		"OpenVPN client status",
		nil,
		nil,
	)
)

// Collector tests network connectivity by querying the IP address location through ip-api.com
type Collector struct {
	connection *measurer.Cached[float64]
}

func NewCollector(httpClient *http.Client, interval time.Duration, logger *slog.Logger) prometheus.Collector {
	const target = "https://clients3.google.com/generate_204"
	return &Collector{
		connection: &measurer.Cached[float64]{
			Interval: interval,
			Do: func(ctx context.Context) (float64, error) {
				attrs := make([]any, 0, 3)
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
				resp, err := httpClient.Do(req)
				attrs = append(attrs, slog.Any("err", err))
				if err == nil {
					_ = resp.Body.Close()
					attrs = append(attrs, slog.Any("status", resp.StatusCode))
				}
				var up float64
				if err == nil && resp.StatusCode == http.StatusNoContent {
					up = 1
				}
				attrs = append(attrs, slog.Any("up", up))
				logger.Debug("connectivity check", attrs...)
				return up, nil
			},
		},
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upMetric
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	up, _ := c.connection.Measure(context.Background())
	ch <- prometheus.MustNewConstMetric(upMetric, prometheus.GaugeValue, up)
}
