package collector_breaker

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestCBCollector(t *testing.T) {
	c := collector{}
	defaultConfiguration.OpenDuration = 500 * time.Millisecond
	cbCollector := New("test", &c, slog.Default())

	t.Run("circuit is closed: collection returns metrics", func(tt *testing.T) {
		ch := make(chan prometheus.Metric)
		go func() { cbCollector.Collect(ch); close(ch) }()
		m := <-ch
		assert.Equal(t, `Desc{fqName: "foo", help: "", constLabels: {}, variableLabels: {}}`, m.Desc().String())
	})

	t.Run("circuit is open: collection returns no metrics", func(tt *testing.T) {
		c.err = errors.New("err")
		ch := make(chan prometheus.Metric)
		go func() { cbCollector.Collect(ch); close(ch) }()
		for m := range ch {
			assert.Contains(t, m.Desc().String(), `Desc{fqName: "circuit_breaker_`)
		}
	})

	t.Run("collection works. circuit eventually closes again", func(t *testing.T) {
		c.err = nil
		assert.Eventually(t, func() bool {
			ch := make(chan prometheus.Metric)
			go func() { cbCollector.Collect(ch); close(ch) }()

			for m := range ch {
				if m.Desc().String() == `Desc{fqName: "foo", help: "", constLabels: {}, variableLabels: {}}` {
					return true
				}
			}
			return false
		}, time.Second, 100*time.Millisecond)
	})
}

var _ Collector = collector{}

type collector struct{ err error }

func (c collector) Describe(_ chan<- *prometheus.Desc) {
}

func (c collector) CollectE(ch chan<- prometheus.Metric) error {
	if c.err != nil {
		return c.err
	}
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc("foo", "", nil, nil), prometheus.CounterValue, 1.0)
	return nil
}
