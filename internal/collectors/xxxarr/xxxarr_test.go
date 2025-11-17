package xxxarr

import (
	"bytes"
	"log/slog"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithToken(t *testing.T) {
	const token = "1234"
	ctx := t.Context()

	f := WithToken(token)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, f(ctx, r))
	assert.Equal(t, token, r.Header.Get("X-API-KEY"))

	f = WithToken("")
	r, _ = http.NewRequest(http.MethodGet, "/", nil)
	assert.Error(t, f(ctx, r))
}

func TestSonarrCollector(t *testing.T) {
	client := fakeClient{
		version: "v1.2.3",
		//health: map[string]int{},
		calendar: []string{
			"foo - S01E01 - 1",
			"foo - S01E01 - 1",
			"foo - S01E02 - 2",
			"foo - S01E03 - 3",
			"foo - S01E04 - 4",
		},
		queue: []QueuedItem{
			{Name: "foo - S01E01 - 1", TotalBytes: 100, DownloadedBytes: 75},
			{Name: "foo - S01E02 - 2", TotalBytes: 100, DownloadedBytes: 50},
		},
		library: Library{Monitored: 3, Unmonitored: 1},
		health:  map[string]int{"foo": 1},
	}
	c, err := NewSonarrCollector("http://localhost:8080", "api-key", http.DefaultClient, slog.New(slog.DiscardHandler))
	require.NoError(t, err)
	c.(*Collector).client = &client

	want := `
# HELP mediamon_xxxarr_calendar Upcoming episodes / movies
# TYPE mediamon_xxxarr_calendar gauge
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E01 - 1",url="http://localhost:8080"} 2
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E02 - 2",url="http://localhost:8080"} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E03 - 3",url="http://localhost:8080"} 1
mediamon_xxxarr_calendar{application="sonarr",title="foo - S01E04 - 4",url="http://localhost:8080"} 1

# HELP mediamon_xxxarr_health Server health
# TYPE mediamon_xxxarr_health gauge
mediamon_xxxarr_health{application="sonarr",type="foo",url="http://localhost:8080"} 1

# HELP mediamon_xxxarr_monitored_count Number of Monitored series / movies
# TYPE mediamon_xxxarr_monitored_count gauge
mediamon_xxxarr_monitored_count{application="sonarr",url="http://localhost:8080"} 3

# HELP mediamon_xxxarr_queued_count Episodes / movies being downloaded
# TYPE mediamon_xxxarr_queued_count gauge
mediamon_xxxarr_queued_count{application="sonarr",url="http://localhost:8080"} 2

# HELP mediamon_xxxarr_queued_downloaded_bytes Downloaded size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_downloaded_bytes gauge
mediamon_xxxarr_queued_downloaded_bytes{application="sonarr",title="foo - S01E01 - 1",url="http://localhost:8080"} 75
mediamon_xxxarr_queued_downloaded_bytes{application="sonarr",title="foo - S01E02 - 2",url="http://localhost:8080"} 50

# HELP mediamon_xxxarr_queued_total_bytes Size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_total_bytes gauge
mediamon_xxxarr_queued_total_bytes{application="sonarr",title="foo - S01E01 - 1",url="http://localhost:8080"} 100
mediamon_xxxarr_queued_total_bytes{application="sonarr",title="foo - S01E02 - 2",url="http://localhost:8080"} 100

# HELP mediamon_xxxarr_unmonitored_count Number of Unmonitored series / movies
# TYPE mediamon_xxxarr_unmonitored_count gauge
mediamon_xxxarr_unmonitored_count{application="sonarr",url="http://localhost:8080"} 1

# HELP mediamon_xxxarr_version Version info
# TYPE mediamon_xxxarr_version gauge
mediamon_xxxarr_version{application="sonarr",url="http://localhost:8080",version="v1.2.3"} 1
`
	assert.NoError(t, testutil.CollectAndCompare(c, bytes.NewBufferString(want)))
}

func TestRadarrCollector(t *testing.T) {
	client := fakeClient{
		version: "v1.2.3",
		//health: map[string]int{},
		calendar: []string{
			"1",
			"1",
			"2",
			"3",
			"4",
		},
		queue: []QueuedItem{
			{Name: "1", TotalBytes: 100, DownloadedBytes: 75},
			{Name: "2", TotalBytes: 100, DownloadedBytes: 50},
		},
		library: Library{Monitored: 3, Unmonitored: 1},
	}
	c, err := NewRadarrCollector("http://localhost:8080", "api-key", http.DefaultClient, slog.New(slog.DiscardHandler))
	require.NoError(t, err)
	c.(*Collector).client = &client

	want := `
# HELP mediamon_xxxarr_calendar Upcoming episodes / movies
# TYPE mediamon_xxxarr_calendar gauge
mediamon_xxxarr_calendar{application="radarr",title="1",url="http://localhost:8080"} 2
mediamon_xxxarr_calendar{application="radarr",title="2",url="http://localhost:8080"} 1
mediamon_xxxarr_calendar{application="radarr",title="3",url="http://localhost:8080"} 1
mediamon_xxxarr_calendar{application="radarr",title="4",url="http://localhost:8080"} 1

# HELP mediamon_xxxarr_monitored_count Number of Monitored series / movies
# TYPE mediamon_xxxarr_monitored_count gauge
mediamon_xxxarr_monitored_count{application="radarr",url="http://localhost:8080"} 3

# HELP mediamon_xxxarr_queued_count Episodes / movies being downloaded
# TYPE mediamon_xxxarr_queued_count gauge
mediamon_xxxarr_queued_count{application="radarr",url="http://localhost:8080"} 2

# HELP mediamon_xxxarr_queued_downloaded_bytes Downloaded size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_downloaded_bytes gauge
mediamon_xxxarr_queued_downloaded_bytes{application="radarr",title="1",url="http://localhost:8080"} 75
mediamon_xxxarr_queued_downloaded_bytes{application="radarr",title="2",url="http://localhost:8080"} 50

# HELP mediamon_xxxarr_queued_total_bytes Size of episode / movie being downloaded in bytes
# TYPE mediamon_xxxarr_queued_total_bytes gauge
mediamon_xxxarr_queued_total_bytes{application="radarr",title="1",url="http://localhost:8080"} 100
mediamon_xxxarr_queued_total_bytes{application="radarr",title="2",url="http://localhost:8080"} 100

# HELP mediamon_xxxarr_unmonitored_count Number of Unmonitored series / movies
# TYPE mediamon_xxxarr_unmonitored_count gauge
mediamon_xxxarr_unmonitored_count{application="radarr",url="http://localhost:8080"} 1

# HELP mediamon_xxxarr_version Version info
# TYPE mediamon_xxxarr_version gauge
mediamon_xxxarr_version{application="radarr",url="http://localhost:8080",version="v1.2.3"} 1
`
	assert.NoError(t, testutil.CollectAndCompare(c, bytes.NewBufferString(want)))
}
