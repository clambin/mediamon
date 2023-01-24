package xxxarr_test

import (
	"bytes"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/collectors/xxxarr/scraper"
	mocks2 "github.com/clambin/mediamon/collectors/xxxarr/scraper/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCollector(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var c *xxxarr.Collector
			switch tt.collector {
			case "sonarr":
				c = xxxarr.NewSonarrCollector("", "")
			case "radarr":
				c = xxxarr.NewRadarrCollector("", "")
			default:
				t.Fatalf("invalid collector type: %s", tt.collector)
			}
			s := mocks2.NewScraper(t)
			c.Scraper = s
			s.On("Scrape", mock.AnythingOfType("*context.emptyCtx")).Return(tt.input, nil)

			r := prometheus.NewPedanticRegistry()
			r.MustRegister(c)

			err := testutil.GatherAndCompare(r, bytes.NewBufferString(tt.output))
			assert.NoError(t, err)

		})
	}
}

type testCase struct {
	name      string
	collector string
	input     scraper.Stats
	output    string
}

var testCases = []testCase{
	{
		name:      "sonarr",
		collector: "sonarr",
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
	{
		name:      "radarr",
		collector: "radarr",
		input: scraper.Stats{
			Version:  "foo",
			Calendar: []string{"1", "2", "3", "4", "5"},
			Queued: []scraper.QueuedFile{
				{Name: "1", TotalBytes: 100, DownloadedBytes: 75},
				{Name: "2", TotalBytes: 100, DownloadedBytes: 50},
				{Name: "2", TotalBytes: 100, DownloadedBytes: 25},
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
mediamon_xxxarr_queued_count{application="radarr",url=""} 3
# HELP mediamon_xxxarr_queued_downloaded_bytes Downloaded size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_downloaded_bytes gauge
mediamon_xxxarr_queued_downloaded_bytes{application="radarr",title="1",url=""} 75
mediamon_xxxarr_queued_downloaded_bytes{application="radarr",title="2",url=""} 75
# HELP mediamon_xxxarr_queued_total_bytes Size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_total_bytes gauge
mediamon_xxxarr_queued_total_bytes{application="radarr",title="1",url=""} 100
mediamon_xxxarr_queued_total_bytes{application="radarr",title="2",url=""} 200
# HELP mediamon_xxxarr_unmonitored_count Number of Unmonitored series / movies
# TYPE mediamon_xxxarr_unmonitored_count gauge
mediamon_xxxarr_unmonitored_count{application="radarr",url=""} 1
# HELP mediamon_xxxarr_version Version info
# TYPE mediamon_xxxarr_version gauge
mediamon_xxxarr_version{application="radarr",url="",version="foo"} 1
`,
	},
}
