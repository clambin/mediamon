package collector_breaker

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestCBCollector(t *testing.T) {
	c := collector{metric: prometheus.NewDesc("foo", "", nil, nil)}
	defaultConfiguration.OpenDuration = 500 * time.Millisecond
	cbCollector := New("test", &c, slog.Default())

	t.Run("circuit is closed: collection returns metrics", func(tt *testing.T) {
		assert.NoError(tt, testutil.CollectAndCompare(
			cbCollector,
			strings.NewReader(`
# HELP foo 
# TYPE foo counter
foo 1
`),
			"foo",
		))
	})

	t.Run("circuit is open: collection returns no metrics", func(tt *testing.T) {
		c.err = errors.New("err")
		assert.Zero(tt, testutil.CollectAndCount(cbCollector, "foo"))
	})

	t.Run("collection works. circuit eventually closes again", func(t *testing.T) {
		c.err = nil
		assert.Eventually(t, func() bool {
			return nil == testutil.CollectAndCompare(
				cbCollector,
				strings.NewReader(`
# HELP foo 
# TYPE foo counter
foo 1
`),
				"foo",
			)
		}, time.Second, 100*time.Millisecond)
	})
}

var _ Collector = collector{}

type collector struct {
	metric *prometheus.Desc
	err    error
}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.metric
}

func (c collector) CollectE(ch chan<- prometheus.Metric) error {
	if c.err == nil {
		ch <- prometheus.MustNewConstMetric(c.metric, prometheus.CounterValue, 1.0)
	}
	return c.err
}
