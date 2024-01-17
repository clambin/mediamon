package xxxarr

import (
	"bytes"
	"context"
	"errors"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollector(t *testing.T) {
	testCases := []struct {
		name      string
		collector string
		setup     func(*mocks.Client, context.Context)
		want      string
	}{
		{
			name:      "sonarr",
			collector: "sonarr",
			setup: func(c *mocks.Client, ctx context.Context) {
				c.EXPECT().GetVersion(ctx).Return("foo", nil)
				c.EXPECT().GetHealth(ctx).Return(nil, nil)
				c.EXPECT().GetCalendar(ctx).Return([]string{
					"foo - S01E01 - 1",
					"foo - S01E02 - 2",
					"foo - S01E03 - 3",
					"foo - S01E04 - 4",
					"foo - S01E05 - 5",
				}, nil)
				c.EXPECT().GetQueue(ctx).Return([]clients.QueuedItem{
					{Name: "foo - S01E01 - 1", TotalBytes: 100, DownloadedBytes: 75},
					{Name: "foo - S01E02 - 2", TotalBytes: 100, DownloadedBytes: 50},
				}, nil)
				c.EXPECT().GetLibrary(ctx).Return(clients.Library{Monitored: 3, Unmonitored: 1}, nil)
			},
			want: `
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
			name:      "sonarr - duplicates",
			collector: "sonarr",
			setup: func(c *mocks.Client, ctx context.Context) {
				c.EXPECT().GetVersion(ctx).Return("foo", nil)
				c.EXPECT().GetHealth(ctx).Return(nil, nil)
				c.EXPECT().GetCalendar(ctx).Return([]string{
					"foo - S01E01 - 1",
					"foo - S01E01 - 1",
					"foo - S01E02 - 2",
					"foo - S01E03 - 3",
					"foo - S01E04 - 4",
				}, nil)
				c.EXPECT().GetQueue(ctx).Return([]clients.QueuedItem{
					{Name: "foo - S01E01 - 1", TotalBytes: 100, DownloadedBytes: 75},
					{Name: "foo - S01E02 - 2", TotalBytes: 100, DownloadedBytes: 50},
				}, nil)
				c.EXPECT().GetLibrary(ctx).Return(clients.Library{Monitored: 3, Unmonitored: 1}, nil)
			},
			want: `
# HELP mediamon_xxxarr_calendar Upcoming episodes / movies
# TYPE mediamon_xxxarr_calendar gauge
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E01 - 1",url=""} 2
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E02 - 2",url=""} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E03 - 3",url=""} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E04 - 4",url=""} 1
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
			name:      "sonarr (error)",
			collector: "sonarr",
			setup: func(c *mocks.Client, ctx context.Context) {
				err := errors.New("sonarr down")
				c.EXPECT().GetVersion(ctx).Return("", err)
				c.EXPECT().GetHealth(ctx).Return(nil, err)
				c.EXPECT().GetCalendar(ctx).Return(nil, err)
				c.EXPECT().GetQueue(ctx).Return(nil, err)
				c.EXPECT().GetLibrary(ctx).Return(clients.Library{}, err)
			},
		},
		{
			name:      "radarr",
			collector: "radarr",
			setup: func(c *mocks.Client, ctx context.Context) {
				c.EXPECT().GetVersion(ctx).Return("foo", nil)
				c.EXPECT().GetHealth(ctx).Return(nil, nil)
				c.EXPECT().GetCalendar(ctx).Return([]string{
					"1",
					"2",
					"3",
					"4",
					"5",
				}, nil)
				c.EXPECT().GetQueue(ctx).Return([]clients.QueuedItem{
					{Name: "1", TotalBytes: 100, DownloadedBytes: 75},
					{Name: "2", TotalBytes: 100, DownloadedBytes: 50},
				}, nil)
				c.EXPECT().GetLibrary(ctx).Return(clients.Library{Monitored: 2, Unmonitored: 1}, nil)

			},
			want: `
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
		{
			name:      "radarr (error)",
			collector: "radarr",
			setup: func(c *mocks.Client, ctx context.Context) {
				err := errors.New("sonarr down")
				c.EXPECT().GetVersion(ctx).Return("", err)
				c.EXPECT().GetHealth(ctx).Return(nil, err)
				c.EXPECT().GetCalendar(ctx).Return(nil, err)
				c.EXPECT().GetQueue(ctx).Return(nil, err)
				c.EXPECT().GetLibrary(ctx).Return(clients.Library{}, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var c *Collector
			switch tt.collector {
			case "sonarr":
				c = NewSonarrCollector("", "")
			case "radarr":
				c = NewRadarrCollector("", "")
			default:
				t.Fatalf("invalid collector type: %s", tt.collector)
			}
			client := mocks.NewClient(t)
			tt.setup(client, context.Background())
			c.client = client

			err := testutil.CollectAndCompare(c, bytes.NewBufferString(tt.want))
			assert.NoError(t, err)
		})
	}
}
