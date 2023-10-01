package scraper_test

import (
	"context"
	"github.com/clambin/mediaclients/xxxarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/scraper"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/scraper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRadarrScraper_Scrape(t *testing.T) {
	c := mocks.NewRadarrGetter(t)
	u := scraper.RadarrScraper{Radarr: c}

	c.EXPECT().GetURL().Return("http://localhost:8080")
	c.EXPECT().GetSystemStatus(mock.Anything).Return(radarrSystemStatus, nil)
	c.EXPECT().GetHealth(mock.Anything).Return(radarrSystemHealth, nil)
	c.EXPECT().GetCalendar(mock.Anything).Return(radarrCalendar, nil)
	c.EXPECT().GetQueue(mock.Anything).Return(radarrQueue, nil)
	c.EXPECT().GetMovies(mock.Anything).Return(radarrMovies, nil)
	for id, entry := range radarrMoviesByID {
		c.EXPECT().GetMovieByID(mock.Anything, id).Return(entry, nil).Once()
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
