package clients

import (
	"context"
	"github.com/clambin/mediaclients/xxxarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRadarr_GetCalendar(t *testing.T) {
	ctx := context.Background()
	r := mocks.NewRadarrClient(t)
	r.EXPECT().GetCalendar(ctx).Return(RadarrCalendar, nil)
	c := Radarr{Client: r}

	calendar, err := c.GetCalendar(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"movie 1", "movie 2", "movie 3"}, calendar)
}

func TestRadarr_GetHealth(t *testing.T) {
	ctx := context.Background()
	r := mocks.NewRadarrClient(t)
	r.EXPECT().GetHealth(ctx).Return(RadarrSystemHealth, nil)
	c := Radarr{Client: r}

	health, err := c.GetHealth(ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"ok": 1, "warning": 1}, health)
}

func TestRadarr_GetLibrary(t *testing.T) {
	ctx := context.Background()
	r := mocks.NewRadarrClient(t)
	r.EXPECT().GetMovies(ctx).Return(RadarrMovies, nil)
	c := Radarr{Client: r}

	library, err := c.GetLibrary(ctx)
	assert.NoError(t, err)
	assert.Equal(t, Library{Monitored: 3, Unmonitored: 1}, library)
}

func TestRadarr_GetQueue(t *testing.T) {
	ctx := context.Background()
	r := mocks.NewRadarrClient(t)
	r.EXPECT().GetQueue(ctx).Return(RadarrQueue, nil)
	c := Radarr{Client: r}

	library, err := c.GetQueue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []QueuedItem{
		{Name: "file1", TotalBytes: 100, DownloadedBytes: 50},
		{Name: "file3", TotalBytes: 100, DownloadedBytes: 100},
		{Name: "file4", TotalBytes: 100, DownloadedBytes: 75},
	}, library)
}

func TestRadarr_GetVersion(t *testing.T) {
	ctx := context.Background()
	r := mocks.NewRadarrClient(t)
	r.EXPECT().GetSystemStatus(ctx).Return(RadarrSystemStatus, nil)
	c := Radarr{Client: r}

	version, err := c.GetVersion(ctx)
	assert.NoError(t, err)
	assert.Equal(t, SonarrSystemStatus.Version, version)
}

var (
	RadarrSystemStatus = xxxarr.RadarrSystemStatusResponse{
		Version: "1.2.3.4444",
	}

	RadarrSystemHealth = []xxxarr.RadarrHealthResponse{
		{
			Type: "ok",
		},
		{
			Type: "warning",
		},
	}

	RadarrCalendar = []xxxarr.RadarrCalendarResponse{
		{Title: "movie 1", Monitored: true, HasFile: false},
		{Title: "movie 2", Monitored: false, HasFile: false},
		{Title: "movie 3", Monitored: true, HasFile: true},
	}

	RadarrQueue = xxxarr.RadarrQueueResponse{
		Page:         1,
		PageSize:     10,
		TotalRecords: 3,
		Records: []xxxarr.RadarrQueueResponseRecord{
			{MovieID: 1, Title: "file1", Status: "downloading", Size: 100, Sizeleft: 50},
			{MovieID: 2, Title: "file3", Status: "downloaded???", Size: 100, Sizeleft: 0},
			{MovieID: 3, Title: "file4", Status: "downloading", Size: 100, Sizeleft: 25},
		},
	}

	RadarrMovies = []xxxarr.RadarrMovieResponse{
		{Title: "movie 1", Monitored: true},
		{Title: "movie 2", Monitored: false},
		{Title: "movie 3", Monitored: true},
		{Title: "movie 5", Monitored: true},
	}

	/*
		RadarrMoviesByID = map[int]xxxarr.RadarrMovieResponse{
			1: {Title: "movie 1"},
			2: {Title: "movie 3"},
			3: {Title: "movie 4"},
		}

	*/
)
