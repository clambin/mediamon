package xxxarr_test

import (
	"fmt"
	"github.com/clambin/go-metrics/tools"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/collectors/xxxarr/scraper"
	mocks2 "github.com/clambin/mediamon/collectors/xxxarr/scraper/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
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
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(testCases["sonarr"].output)))
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
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(testCases["radarr"].output)))

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
	assert.Equal(t, "mediamon_error", tools.MetricName(metric))
	assert.Contains(t, metric.Desc().String(), "Error getting sonarr metrics")

}

type testCase struct {
	input  scraper.Stats
	output string
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
		output: `
# HELP mediamon_xxxarr_calendar Upcoming episodes / movies
# TYPE mediamon_xxxarr_calendar gauge
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E01 - 1",url=""} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E02 - 2",url=""} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E03 - 3",url=""} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E04 - 4",url=""} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E05 - 5",url=""} 1
# HELP mediamon_xxxarr_monitored_count Number of Monitored series / movies
# TYPE mediamon_xxxarr_monitored_count gauge
mediamon_xxxarr_monitored_count{application="sonarr",url=""} 3
# HELP mediamon_xxxarr_queued_count Episodes / movies being downloaded
# TYPE mediamon_xxxarr_queued_count gauge
mediamon_xxxarr_queued_count{application="sonarr",url=""} 2
# HELP mediamon_xxxarr_queued_downloaded_bytes Downloaded size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_downloaded_bytes gauge
mediamon_xxxarr_queued_downloaded_bytes{application="sonarr",title="foo - S01E01 - 1",url=""} 75
mediamon_xxxarr_queued_downloaded_bytes{application="sonarr",title="foo - S01E02 - 2",url=""} 50
# HELP mediamon_xxxarr_queued_total_bytes Size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_total_bytes gauge
mediamon_xxxarr_queued_total_bytes{application="sonarr",title="foo - S01E01 - 1",url=""} 100
mediamon_xxxarr_queued_total_bytes{application="sonarr",title="foo - S01E02 - 2",url=""} 100
# HELP mediamon_xxxarr_unmonitored_count Number of Unmonitored series / movies
# TYPE mediamon_xxxarr_unmonitored_count gauge
mediamon_xxxarr_unmonitored_count{application="sonarr",url=""} 1
# HELP mediamon_xxxarr_version Version info
# TYPE mediamon_xxxarr_version gauge
mediamon_xxxarr_version{application="sonarr",url="",version="foo"} 1
`,
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
		output: `
# HELP mediamon_xxxarr_calendar Upcoming episodes / movies
# TYPE mediamon_xxxarr_calendar gauge
mediamon_xxxarr_calendar{application="radarr",title="1",url=""} 1
mediamon_xxxarr_calendar{application="radarr",title="2",url=""} 1
mediamon_xxxarr_calendar{application="radarr",title="3",url=""} 1
mediamon_xxxarr_calendar{application="radarr",title="4",url=""} 1
mediamon_xxxarr_calendar{application="radarr",title="5",url=""} 1
# HELP mediamon_xxxarr_monitored_count Number of Monitored series / movies
# TYPE mediamon_xxxarr_monitored_count gauge
mediamon_xxxarr_monitored_count{application="radarr",url=""} 2
# HELP mediamon_xxxarr_queued_count Episodes / movies being downloaded
# TYPE mediamon_xxxarr_queued_count gauge
mediamon_xxxarr_queued_count{application="radarr",url=""} 2
# HELP mediamon_xxxarr_queued_downloaded_bytes Downloaded size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_downloaded_bytes gauge
mediamon_xxxarr_queued_downloaded_bytes{application="radarr",title="1",url=""} 75
mediamon_xxxarr_queued_downloaded_bytes{application="radarr",title="2",url=""} 50
# HELP mediamon_xxxarr_queued_total_bytes Size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_total_bytes gauge
mediamon_xxxarr_queued_total_bytes{application="radarr",title="1",url=""} 100
mediamon_xxxarr_queued_total_bytes{application="radarr",title="2",url=""} 100
# HELP mediamon_xxxarr_unmonitored_count Number of Unmonitored series / movies
# TYPE mediamon_xxxarr_unmonitored_count gauge
mediamon_xxxarr_unmonitored_count{application="radarr",url=""} 1
# HELP mediamon_xxxarr_version Version info
# TYPE mediamon_xxxarr_version gauge
mediamon_xxxarr_version{application="radarr",url="",version="foo"} 1
`,
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
