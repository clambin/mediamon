// Metrics is a utility package to facilitate unit tests by allowing
// the test to read back to the value set in a Prometheus metrics
package metrics

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

var (
	metrics = make(map[string]interface{})
)

// Gauge returns a new Prometheus GaugeVec, created through promauto
func NewGaugeVec(opts prometheus.GaugeOpts, labels []string) *prometheus.GaugeVec {
	metric := promauto.NewGaugeVec(opts, labels)
	metrics[opts.Name] = metric

	return metric
}

// Gauge returns a new Prometheus Gauge, created through promauto
func NewGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	metric := promauto.NewGauge(opts)
	metrics[opts.Name] = metric

	return metric
}

// LoadValue gets the last value reported so unit tests can verify the correct value was reported
func LoadValue(metricName string, labels ...string) (float64, error) {
	log.Debugf("%s(%s)", metricName, labels)
	if metric, ok := metrics[metricName]; ok {
		var m = dto.Metric{}
		switch metricType := metric.(type) {
		case prometheus.Gauge:
			gauge := metric.(prometheus.Gauge)
			_ = gauge.Write(&m)
			return m.Gauge.GetValue(), nil
		case *prometheus.GaugeVec:
			gaugevec := metric.(*prometheus.GaugeVec)
			log.Debug(gaugevec)
			_ = gaugevec.WithLabelValues(labels...).Write(&m)
			return m.Gauge.GetValue(), nil
		default:
			return 0, errors.New(fmt.Sprintf("invalid type for metric %s: %v", metricName, metricType))
		}
	}
	return 0, errors.New("could not find " + metricName)
}
