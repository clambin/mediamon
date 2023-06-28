package scraper_test

import (
	"context"
	"github.com/clambin/mediamon/v2/collectors/xxxarr/scraper"
	"github.com/clambin/mediamon/v2/collectors/xxxarr/scraper/mocks"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/xxxarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRadarrScraper_Scrape(t *testing.T) {
	c := mocks.NewRadarrAPI(t)
	u := scraper.RadarrScraper{Client: c}

	c.On("GetURL").Return("http://localhost:8080")
	c.On("GetSystemStatus", mock.AnythingOfType("*context.emptyCtx")).Return(radarrSystemStatus, nil)
	c.On("GetHealth", mock.AnythingOfType("*context.emptyCtx")).Return(radarrSystemHealth, nil)
	c.On("GetCalendar", mock.AnythingOfType("*context.emptyCtx")).Return(radarrCalendar, nil)
	c.On("GetQueue", mock.AnythingOfType("*context.emptyCtx")).Return(radarrQueue, nil)
	c.On("GetMovies", mock.AnythingOfType("*context.emptyCtx")).Return(radarrMovies, nil)
	for id, entry := range radarrMoviesByID {
		c.On("GetMovieByID", mock.AnythingOfType("*context.emptyCtx"), id).Return(entry, nil).Once()

	}

	stats, err := u.Scrape(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8080", stats.URL)
	assert.Equal(t, 1, stats.Health["ok"])
	assert.Equal(t, 1, stats.Health["warning"])
	assert.Equal(t, "1.2.3.4444", stats.Version)
	assert.Equal(t, []string{"movie 1", "movie 2"}, stats.Calendar)
	assert.Equal(t, []scraper.QueuedFile{
		{Name: "movie 1", TotalBytes: 100, DownloadedBytes: 50},
		{Name: "movie 3", TotalBytes: 100, DownloadedBytes: 100},
		{Name: "movie 4", TotalBytes: 100, DownloadedBytes: 75},
	}, stats.Queued)
	assert.Equal(t, 3, stats.Monitored)
	assert.Equal(t, 1, stats.Unmonitored)
}

var (
	radarrSystemStatus = xxxarr.RadarrSystemStatusResponse{
		Version: "1.2.3.4444",
	}

	radarrSystemHealth = []xxxarr.RadarrHealthResponse{
		{
			Type: "ok",
		},
		{
			Type: "warning",
		},
	}

	radarrCalendar = []xxxarr.RadarrCalendarResponse{
		{Title: "movie 1", Monitored: true, HasFile: false},
		{Title: "movie 2", Monitored: false, HasFile: false},
		{Title: "movie 3", Monitored: true, HasFile: true},
	}

	radarrQueue = xxxarr.RadarrQueueResponse{
		Page:         1,
		PageSize:     10,
		TotalRecords: 3,
		Records: []xxxarr.RadarrQueueResponseRecord{
			{MovieID: 1, Title: "file1", Status: "downloading", Size: 100, Sizeleft: 50},
			{MovieID: 2, Title: "file3", Status: "downloaded???", Size: 100, Sizeleft: 0},
			{MovieID: 3, Title: "file4", Status: "downloading", Size: 100, Sizeleft: 25},
		},
	}

	radarrMovies = []xxxarr.RadarrMovieResponse{
		{Title: "movie 1", Monitored: true},
		{Title: "movie 2", Monitored: false},
		{Title: "movie 3", Monitored: true},
		{Title: "movie 5", Monitored: true},
	}

	radarrMoviesByID = map[int]xxxarr.RadarrMovieResponse{
		1: {Title: "movie 1"},
		2: {Title: "movie 3"},
		3: {Title: "movie 4"},
	}
)
