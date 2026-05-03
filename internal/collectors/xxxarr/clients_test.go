package xxxarr

import (
	"net/http"
	"testing"

	"github.com/clambin/mediaclients/radarr"
	"github.com/clambin/mediaclients/sonarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRadarrClient(t *testing.T) {
	client := fakeRadarrClient{
		systemStatus: &radarr.GetApiV3SystemStatusResponse{JSON200: &radarr.SystemResource{Version: new("v1.2.3")}},
		health: &radarr.GetApiV3HealthResponse{JSON200: &[]radarr.HealthResource{
			{Type: new(radarr.HealthCheckResult("foo")), Message: new("bar")},
		}},
		calendar: &radarr.GetApiV3CalendarResponse{JSON200: &[]radarr.MovieResource{{Title: new("some movie")}}},
		queue: &radarr.GetApiV3QueueResponse{JSON200: &radarr.QueueResourcePagingResource{
			Page:         new(int32(1)),
			PageSize:     new(int32(100)),
			Records:      &[]radarr.QueueResource{{Size: new(100.0), Sizeleft: new(40.0), Title: new("some other movie")}},
			TotalRecords: new(int32(1)),
		}},
		movies: &radarr.GetApiV3MovieResponse{JSON200: &[]radarr.MovieResource{
			{Monitored: new(true), Title: new("some movie")},
			{Monitored: new(false), Title: new("some other movie")},
			{Monitored: new(true), Title: new("some other other movie")},
		}},
	}
	c, _ := NewRadarrClient("http://localhost:1234", "api-key", http.DefaultClient)
	c.Client = &client

	ctx := t.Context()
	version, err := c.GetVersion(ctx)
	require.NoError(t, err)
	assert.Equal(t, "v1.2.3", version)

	health, err := c.GetHealth(ctx)
	require.NoError(t, err)
	assert.Equal(t, map[string]int{"foo": 1}, health)

	calendar, err := c.GetCalendar(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"some movie"}, calendar)

	queue, err := c.GetQueue(ctx)
	require.NoError(t, err)
	assert.Equal(t, []QueuedItem{{Name: "some other movie", TotalBytes: 100, DownloadedBytes: 60}}, queue)

	library, err := c.GetLibrary(ctx)
	require.NoError(t, err)
	assert.Equal(t, Library{Monitored: 2, Unmonitored: 1}, library)
}

func TestSonarrClient(t *testing.T) {
	client := fakeSonarrClient{
		systemStatus: &sonarr.GetApiV3SystemStatusResponse{JSON200: &sonarr.SystemResource{Version: new("v1.2.3")}},
		health: &sonarr.GetApiV3HealthResponse{JSON200: &[]sonarr.HealthResource{
			{Type: new(sonarr.HealthCheckResult("foo")), Message: new("bar")},
		}},
		calendar: &sonarr.GetApiV3CalendarResponse{JSON200: &[]sonarr.EpisodeResource{{
			Title:         new("some episode"),
			SeasonNumber:  new(int32(1)),
			EpisodeNumber: new(int32(12)),
			Series:        &sonarr.SeriesResource{Title: new("some series")}},
		}},
		queue: &sonarr.GetApiV3QueueResponse{JSON200: &sonarr.QueueResourcePagingResource{
			Page:     new(int32(1)),
			PageSize: new(int32(100)),
			Records: &[]sonarr.QueueResource{{
				Size:         new(100.0),
				Sizeleft:     new(40.0),
				Title:        new("some other episode"),
				Series:       &sonarr.SeriesResource{Title: new("some other series")},
				SeasonNumber: new(int32(1)),
				Episode:      &sonarr.EpisodeResource{EpisodeNumber: new(int32(12))},
			}},
			TotalRecords: new(int32(1)),
		}},
		series: &sonarr.GetApiV3SeriesResponse{JSON200: &[]sonarr.SeriesResource{
			{Monitored: new(true), Title: new("some series")},
			{Monitored: new(false), Title: new("some other series")},
			{Monitored: new(true), Title: new("some other other series")},
		}},
	}
	c, _ := NewSonarrClient("http://localhost:1234", "api-key", http.DefaultClient)
	c.Client = &client

	ctx := t.Context()
	version, err := c.GetVersion(ctx)
	require.NoError(t, err)
	assert.Equal(t, "v1.2.3", version)

	health, err := c.GetHealth(ctx)
	require.NoError(t, err)
	assert.Equal(t, map[string]int{"foo": 1}, health)

	calendar, err := c.GetCalendar(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"some series - S01E12 - some episode"}, calendar)

	queue, err := c.GetQueue(ctx)
	require.NoError(t, err)
	assert.Equal(t, []QueuedItem{{Name: "some other series - S01E12 - some other episode", TotalBytes: 100, DownloadedBytes: 60}}, queue)

	library, err := c.GetLibrary(ctx)
	require.NoError(t, err)
	assert.Equal(t, Library{Monitored: 2, Unmonitored: 1}, library)
}
