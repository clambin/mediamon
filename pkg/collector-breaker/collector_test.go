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
	cbCollector := New(&c, 3, 500*time.Millisecond, 1, slog.Default())

	cbCollector.Describe(make(chan *prometheus.Desc))

	// circuit is closed: collection returns metrics
	ch := make(chan prometheus.Metric)
	go func() { cbCollector.Collect(ch); close(ch) }()
	m := <-ch
	assert.Equal(t, `Desc{fqName: "foo", help: "", constLabels: {}, variableLabels: {}}`, m.Desc().String())

	// circuit is open: collection returns no metrics
	c.err = errors.New("err")
	ch = make(chan prometheus.Metric)
	go func() { cbCollector.Collect(ch); close(ch) }()
	_, ok := <-ch
	assert.False(t, ok)

	// collection works. circuit eventually closes again
	c.err = nil
	assert.Eventually(t, func() bool {
		ch = make(chan prometheus.Metric)
		go func() { cbCollector.Collect(ch); close(ch) }()

		m = <-ch
		return m.Desc().String() == `Desc{fqName: "foo", help: "", constLabels: {}, variableLabels: {}}`
	}, time.Second, 100*time.Millisecond)
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
