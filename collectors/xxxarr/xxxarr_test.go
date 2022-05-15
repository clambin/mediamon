package xxxarr_test

import (
	"fmt"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/collectors/xxxarr/scraper"
	mocks2 "github.com/clambin/mediamon/collectors/xxxarr/scraper/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestSonarrCollector_Describe(t *testing.T) {
	c := xxxarr.NewSonarrCollector("http://localhost:8888", "")
	testCollectorDescribe(t, c, `constLabels: {application="sonarr"}`)
}

func TestSonarrCollector_Collect(t *testing.T) {
	c := xxxarr.NewSonarrCollector("", "")
	s := &mocks2.Scraper{}
	c.(*xxxarr.Collector).Scraper = s
	s.On("Scrape").Return(testCases["sonarr"].input, nil)
	testCollectorCollect(t, c, testCases["sonarr"].output)
	s.AssertExpectations(t)
}

func TestRadarrCollector_Describe(t *testing.T) {
	c := xxxarr.NewRadarrCollector("http://localhost:8888", "")
	testCollectorDescribe(t, c, `constLabels: {application="radarr"}`)
}

func TestRadarrCollector_Collect(t *testing.T) {
	c := xxxarr.NewRadarrCollector("", "")
	s := &mocks2.Scraper{}
	c.(*xxxarr.Collector).Scraper = s
	s.On("Scrape").Return(testCases["radarr"].input, nil)
	testCollectorCollect(t, c, testCases["radarr"].output)
	s.AssertExpectations(t)
}

func TestCollector_Failure(t *testing.T) {
	c := xxxarr.NewSonarrCollector("", "")
	s := &mocks2.Scraper{}
	c.(*xxxarr.Collector).Scraper = s
	s.On("Scrape").Return(scraper.Stats{}, fmt.Errorf("failure"))

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	metric := <-ch
	assert.Equal(t, "mediamon_error", metrics.MetricName(metric))
	assert.Contains(t, metric.Desc().String(), "Error getting sonarr metrics")

}

type testCase struct {
	input  scraper.Stats
	output expectedOutput
}

type expectedOutput struct {
	application string
	version     string
	calendar    []string
	queued      []expectedQueue
	monitored   float64
	unmonitored float64
}

type expectedQueue struct {
	name       string
	size       float64
	downloaded float64
}

var testCases = map[string]testCase{
	"sonarr": {
		input: scraper.Stats{
			URL:      "",
			Version:  "foo",
			Calendar: []string{"foo - S01E01 - 1", "foo - S01E02 - 2", "foo - S01E03 - 3", "foo - S01E04 - 4", "foo - S01E05 - 5"},
			Queued: []scraper.QueuedFile{
				{Name: "foo - S01E01 - 1", TotalBytes: 100, DownloadedBytes: 75},
				{Name: "foo - S01E02 - 2", TotalBytes: 100, DownloadedBytes: 50},
			},
			Monitored:   3,
			Unmonitored: 1,
		},
		output: expectedOutput{
			application: "sonarr",
			version:     "foo",
			calendar:    []string{"foo - S01E01 - 1", "foo - S01E02 - 2", "foo - S01E03 - 3", "foo - S01E04 - 4", "foo - S01E05 - 5"},
			queued: []expectedQueue{
				{name: "foo - S01E01 - 1", size: 100, downloaded: 75},
				{name: "foo - S01E02 - 2", size: 100, downloaded: 50},
			},
			monitored:   3,
			unmonitored: 1,
		},
	},
	"radarr": {
		input: scraper.Stats{
			Version:  "foo",
			Calendar: []string{"1", "2", "3", "4", "5"},
			Queued: []scraper.QueuedFile{
				{Name: "1", TotalBytes: 100, DownloadedBytes: 75},
				{Name: "2", TotalBytes: 100, DownloadedBytes: 50},
			},
			Monitored:   2,
			Unmonitored: 1,
		},
		output: expectedOutput{
			application: "radarr",
			version:     "foo",
			calendar:    []string{"1", "2", "3", "4", "5"},
			queued: []expectedQueue{
				{name: "1", size: 100, downloaded: 75},
				{name: "2", size: 100, downloaded: 50},
			},
			monitored:   2,
			unmonitored: 1,
		},
	},
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

func testCollectorCollect(t *testing.T, collector prometheus.Collector, expected expectedOutput) {
	ch := make(chan prometheus.Metric)
	go collector.Collect(ch)

	// version
	metric := <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, "foo", metrics.MetricLabel(metric, "version"))

	// calendar
	for _, title := range expected.calendar {
		metric = <-ch
		assert.Equal(t, 1.0, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, expected.application, metrics.MetricLabel(metric, "application"))
		assert.Equal(t, title, metrics.MetricLabel(metric, "title"))

	}

	metric = <-ch
	assert.Equal(t, float64(len(expected.queued)), metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, expected.application, metrics.MetricLabel(metric, "application"))

	for _, entry := range expected.queued {
		metric = <-ch
		assert.Equal(t, entry.size, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, expected.application, metrics.MetricLabel(metric, "application"))
		assert.Equal(t, entry.name, metrics.MetricLabel(metric, "title"))

		metric = <-ch
		assert.Equal(t, entry.downloaded, metrics.MetricValue(metric).GetGauge().GetValue())
		assert.Equal(t, expected.application, metrics.MetricLabel(metric, "application"))
		assert.Equal(t, entry.name, metrics.MetricLabel(metric, "title"))
	}

	// monitored
	metric = <-ch
	assert.Equal(t, expected.monitored, metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, expected.application, metrics.MetricLabel(metric, "application"))

	// unmonitored
	metric = <-ch
	assert.Equal(t, expected.unmonitored, metrics.MetricValue(metric).GetGauge().GetValue())
	assert.Equal(t, expected.application, metrics.MetricLabel(metric, "application"))

}
