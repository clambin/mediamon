package clients

import (
	"context"
	"github.com/clambin/mediaclients/xxxarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSonarr_GetCalendar(t *testing.T) {
	ctx := context.Background()
	s := mocks.NewSonarrClient(t)
	s.EXPECT().GetCalendar(ctx).Return(SonarrCalendar, nil)
	c := Sonarr{Client: s}

	calendar, err := c.GetCalendar(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar", "snafu", "ufans"}, calendar)
}

func TestSonarr_GetHealth(t *testing.T) {
	ctx := context.Background()
	s := mocks.NewSonarrClient(t)
	s.EXPECT().GetHealth(ctx).Return(SonarrSystemHealth, nil)
	c := Sonarr{Client: s}

	health, err := c.GetHealth(ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"error": 1, "warning": 1}, health)
}

func TestSonarr_GetLibrary(t *testing.T) {
	ctx := context.Background()
	s := mocks.NewSonarrClient(t)
	s.EXPECT().GetSeries(ctx).Return(SonarrSeries, nil)
	c := Sonarr{Client: s}

	library, err := c.GetLibrary(ctx)
	assert.NoError(t, err)
	assert.Equal(t, Library{Monitored: 3, Unmonitored: 1}, library)
}

func TestSonarr_GetQueue(t *testing.T) {
	ctx := context.Background()
	s := mocks.NewSonarrClient(t)
	s.EXPECT().GetQueue(ctx).Return(SonarrQueue, nil)
	for key, value := range SonarrEpisodes {
		s.EXPECT().GetEpisodeByID(ctx, key).Return(value, nil)
	}
	c := Sonarr{Client: s}

	library, err := c.GetQueue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []QueuedItem{
		{Name: "series - S01E01 - Pilot", TotalBytes: 100, DownloadedBytes: 50},
		{Name: "series - S01E02 - Seconds", TotalBytes: 100, DownloadedBytes: 100},
		{Name: "series - S01E03 - End", TotalBytes: 100, DownloadedBytes: 75},
	}, library)
}

func TestSonarr_GetVersion(t *testing.T) {
	ctx := context.Background()
	s := mocks.NewSonarrClient(t)
	s.EXPECT().GetSystemStatus(ctx).Return(SonarrSystemStatus, nil)
	c := Sonarr{Client: s}

	version, err := c.GetVersion(ctx)
	assert.NoError(t, err)
	assert.Equal(t, SonarrSystemStatus.Version, version)
}

var (
	SonarrSystemStatus = xxxarr.SonarrSystemStatusResponse{
		Version: "1.2.3.4444",
	}

	SonarrSystemHealth = []xxxarr.SonarrHealthResponse{
		{Type: "warning"},
		{Type: "error"},
	}

	SonarrCalendar = []xxxarr.SonarrCalendarResponse{
		{SeriesID: 11, SeasonNumber: 1, EpisodeNumber: 1, Title: "foo", Monitored: true, HasFile: true},
		{SeriesID: 11, SeasonNumber: 1, EpisodeNumber: 2, Title: "bar", Monitored: true, HasFile: false},
		{SeriesID: 11, SeasonNumber: 1, EpisodeNumber: 3, Title: "snafu", Monitored: true, HasFile: false},
		{SeriesID: 12, SeasonNumber: 2, EpisodeNumber: 1, Title: "ufans", Monitored: false, HasFile: true},
	}

	/*
		SonarrSeriesByID11 = xxxarr.SonarrSeriesResponse{
			Title: "Series 11",
		}
	*/

	SonarrQueue = xxxarr.SonarrQueueResponse{
		Page:         1,
		PageSize:     10,
		TotalRecords: 3,
		Records: []xxxarr.SonarrQueueResponseRecord{
			{Title: "file1", Status: "downloading", EpisodeID: 1, Size: 100, Sizeleft: 50},
			{Title: "file2", Status: "downloaded???", EpisodeID: 2, Size: 100, Sizeleft: 0},
			{Title: "file3", Status: "downloading", EpisodeID: 3, Size: 100, Sizeleft: 25},
		},
	}

	SonarrEpisodes = map[int]xxxarr.SonarrEpisodeResponse{
		1: {Title: "Pilot", SeasonNumber: 1, EpisodeNumber: 1, Series: xxxarr.SonarrEpisodeResponseSeries{Title: "series"}},
		2: {Title: "Seconds", SeasonNumber: 1, EpisodeNumber: 2, Series: xxxarr.SonarrEpisodeResponseSeries{Title: "series"}},
		3: {Title: "End", SeasonNumber: 1, EpisodeNumber: 3, Series: xxxarr.SonarrEpisodeResponseSeries{Title: "series"}},
	}

	SonarrSeries = []xxxarr.SonarrSeriesResponse{
		{Title: "movie 1", Monitored: true},
		{Title: "movie 2", Monitored: false},
		{Title: "movie 3", Monitored: true},
		{Title: "movie 5", Monitored: true},
	}
)
