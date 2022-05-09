package xxxarr_test

import (
	"fmt"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/collectors/xxxarr/updater"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestSonarrCollector_Describe(t *testing.T) {
	c := xxxarr.NewSonarrCollector("http://localhost:8888", "", 5*time.Minute)
	testCollectorDescribe(t, c, `constLabels: {application="sonarr"}`)
}

func TestSonarrCollector_Collect(t *testing.T) {
	s := &server{application: "sonarr"}

	c := xxxarr.NewSonarrCollector("", "", 5*time.Minute)
	c.(*xxxarr.Collector).Cache.Updater = s.GetStats
	testCollectorCollect(t, c, "sonarr")
}

func TestRadarrCollector_Describe(t *testing.T) {
	c := xxxarr.NewRadarrCollector("http://localhost:8888", "", 5*time.Minute)
	testCollectorDescribe(t, c, `constLabels: {application="radarr"}`)
}

func TestRadarrCollector_Collect(t *testing.T) {
	s := &server{application: "sonarr"}

	c := xxxarr.NewRadarrCollector("", "", 5*time.Minute)
	c.(*xxxarr.Collector).Cache.Updater = s.GetStats

	testCollectorCollect(t, c, "radarr")
}

func testCollectorDescribe(t *testing.T, collector prometheus.Collector, labelString string) {
	ch := make(chan *prometheus.Desc)
	go collector.Describe(ch)

	metricsNames := []string{
		"mediamon_xxxarr_version",
		"mediamon_xxxarr_calendar",
		"mediamon_xxxarr_queued_count",
		"mediamon_xxxarr_queued_total_bytes",
		"mediamon_xxxarr_queued_downloaded_bytes",
		"mediamon_xxxarr_monitored_count",
		"mediamon_xxxarr_unmonitored_count",
	}

	r, _ := regexp.Compile(`fqName: "([a-z_]+)",`)

	for range metricsNames {
		metric := <-ch
		metricAsString := metric.String()

		name := r.FindStringSubmatch(metricAsString)
		assert.Len(t, name, 2)
		assert.Contains(t, metricsNames, name[1])

		assert.Contains(t, metricAsString, labelString)
	}
}

func testCollectorCollect(t *testing.T, collector prometheus.Collector, application string) {
	ch := make(chan prometheus.Metric)
	go collector.Collect(ch)

	metric := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, "foo", metrics.MetricLabel(metric, "version"))

	for _, title := range []string{"1", "2", "3", "4", "5"} {
		metric = <-ch
		assert.Equal(t, 1.0, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, application, metrics.MetricLabel(metric, "application"))
		assert.Equal(t, title, metrics.MetricLabel(metric, "title"))

	}

	expectedQueued := []struct {
		name       string
		size       float64
		downloaded float64
	}{
		{name: "1", size: 100, downloaded: 75},
		{name: "2", size: 100, downloaded: 100},
	}

	metric = <-ch
	assert.Equal(t, float64(len(expectedQueued)), metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, application, metrics.MetricLabel(metric, "application"))

	for _, entry := range expectedQueued {
		metric = <-ch
		assert.Equal(t, entry.size, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, application, metrics.MetricLabel(metric, "application"))
		assert.Equal(t, entry.name, metrics.MetricLabel(metric, "title"))

		metric = <-ch
		assert.Equal(t, entry.downloaded, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, application, metrics.MetricLabel(metric, "application"))
		assert.Equal(t, entry.name, metrics.MetricLabel(metric, "title"))
	}

	for _, value := range []float64{10, 3} {
		metric = <-ch
		assert.Equal(t, value, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, application, metrics.MetricLabel(metric, "application"))
	}
}

type server struct {
	application string
	failing     bool
}

func (s server) GetStats() (stats updater.Stats, err error) {
	if s.failing {
		return stats, fmt.Errorf("failing")
	}
	stats = updater.Stats{
		URL:      "https://localhost:4321",
		Version:  "foo",
		Calendar: []string{"1", "2", "3", "4", "5"},
		Queued: []updater.QueuedFile{
			{Name: "1", TotalBytes: 100, DownloadedBytes: 75},
			{Name: "2", TotalBytes: 100, DownloadedBytes: 100},
		},
		Monitored:   10,
		Unmonitored: 3,
	}
	return
}
