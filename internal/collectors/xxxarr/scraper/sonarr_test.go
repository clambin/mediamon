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

func TestSonarrScraper_Scrape(t *testing.T) {
	c := mocks.NewSonarrGetter(t)
	u := scraper.SonarrScraper{Sonarr: c}

	c.EXPECT().GetURL().Return("http://localhost:8080")
	c.EXPECT().GetSystemStatus(mock.Anything).Return(sonarrSystemStatus, nil)
	c.EXPECT().GetHealth(mock.Anything).Return(sonarrSystemHealth, nil)
	c.EXPECT().GetCalendar(mock.Anything).Return(sonarrCalendar, nil)
	c.EXPECT().GetQueue(mock.Anything).Return(sonarrQueue, nil)
	c.EXPECT().GetSeries(mock.Anything).Return(sonarrSeries, nil)
	c.EXPECT().GetSeriesByID(mock.Anything, 11).Return(sonarrSeriesByID11, nil)
	for id, entry := range sonarrEpisodes {
		c.EXPECT().GetEpisodeByID(mock.Anything, id).Return(entry, nil)
	}

	stats, err := u.Scrape(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8080", stats.URL)
	assert.Equal(t, 0, stats.Health["ok"])
	assert.Equal(t, 1, stats.Health["warning"])
	assert.Equal(t, 1, stats.Health["error"])
	assert.Equal(t, "1.2.3.4444", stats.Version)
	assert.Equal(t, []string{"Series 11 - S01E02 - bar", "Series 11 - S01E03 - snafu"}, stats.Calendar)
	assert.Equal(t, []scraper.QueuedFile{
		{Name: "series - S01E01 - Pilot", TotalBytes: 100, DownloadedBytes: 50},
		{Name: "series - S01E02 - Seconds", TotalBytes: 100, DownloadedBytes: 100},
		{Name: "series - S01E03 - End", TotalBytes: 100, DownloadedBytes: 75},
	}, stats.Queued)
	assert.Equal(t, 3, stats.Monitored)
	assert.Equal(t, 1, stats.Unmonitored)
}

var (
	sonarrSystemStatus = xxxarr.SonarrSystemStatusResponse{
		Version: "1.2.3.4444",
	}

	sonarrSystemHealth = []xxxarr.SonarrHealthResponse{
		{Type: "warning"},
		{Type: "error"},
	}

	sonarrCalendar = []xxxarr.SonarrCalendarResponse{
		{SeriesID: 11, SeasonNumber: 1, EpisodeNumber: 1, Title: "foo", Monitored: true, HasFile: true},
		{SeriesID: 11, SeasonNumber: 1, EpisodeNumber: 2, Title: "bar", Monitored: true, HasFile: false},
		{SeriesID: 11, SeasonNumber: 1, EpisodeNumber: 3, Title: "snafu", Monitored: true, HasFile: false},
		{SeriesID: 12, SeasonNumber: 2, EpisodeNumber: 1, Title: "ufans", Monitored: false, HasFile: true},
	}

	sonarrSeriesByID11 = xxxarr.SonarrSeriesResponse{
		Title: "Series 11",
	}

	sonarrQueue = xxxarr.SonarrQueueResponse{
		Page:         1,
		PageSize:     10,
		TotalRecords: 3,
		Records: []xxxarr.SonarrQueueResponseRecord{
			{Title: "file1", Status: "downloading", EpisodeID: 1, Size: 100, Sizeleft: 50},
			{Title: "file2", Status: "downloaded???", EpisodeID: 2, Size: 100, Sizeleft: 0},
			{Title: "file3", Status: "downloading", EpisodeID: 3, Size: 100, Sizeleft: 25},
		},
	}

	sonarrEpisodes = map[int]xxxarr.SonarrEpisodeResponse{
		1: {Title: "Pilot", SeasonNumber: 1, EpisodeNumber: 1, Series: xxxarr.SonarrEpisodeResponseSeries{Title: "series"}},
		2: {Title: "Seconds", SeasonNumber: 1, EpisodeNumber: 2, Series: xxxarr.SonarrEpisodeResponseSeries{Title: "series"}},
		3: {Title: "End", SeasonNumber: 1, EpisodeNumber: 3, Series: xxxarr.SonarrEpisodeResponseSeries{Title: "series"}},
	}

	sonarrSeries = []xxxarr.SonarrSeriesResponse{
		{Title: "movie 1", Monitored: true},
		{Title: "movie 2", Monitored: false},
		{Title: "movie 3", Monitored: true},
		{Title: "movie 5", Monitored: true},
	}
)
