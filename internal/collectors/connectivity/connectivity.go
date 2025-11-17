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
				logger.Debug("testing connectivity")
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
				resp, err := httpClient.Do(req)
				if err == nil {
					_ = resp.Body.Close()
				}
				logger.Debug("connectivity test result", "err", err)
				if err == nil {
					logger.Debug("connectivity test response", "status", resp.StatusCode)
				}
				if err == nil && resp.StatusCode == http.StatusNoContent {
					return 1, nil
				}
				return 0, nil
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
