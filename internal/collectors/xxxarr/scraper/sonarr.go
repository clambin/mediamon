package scraper

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/xxxarr"
)

// SonarrAPI contains all supported Sonarr APIs
//
//go:generate mockery --name SonarrAPI --with-expecter=true
type SonarrAPI interface {
	GetURL() string
	GetSystemStatus(ctx context.Context) (response xxxarr.SonarrSystemStatusResponse, err error)
	GetHealth(ctx context.Context) (response []xxxarr.SonarrHealthResponse, err error)
	GetCalendar(ctx context.Context) (response []xxxarr.SonarrCalendarResponse, err error)
	GetQueuePage(ctx context.Context, pageNr int) (response xxxarr.SonarrQueueResponse, err error)
	GetQueue(ctx context.Context) (response xxxarr.SonarrQueueResponse, err error)
	GetSeries(ctx context.Context) (response []xxxarr.SonarrSeriesResponse, err error)
	GetSeriesByID(ctx context.Context, seriesID int) (response xxxarr.SonarrSeriesResponse, err error)
	GetEpisodeByID(ctx context.Context, episodeID int) (response xxxarr.SonarrEpisodeResponse, err error)
}

// SonarrScraper collects Stats from a Sonarr instance
type SonarrScraper struct {
	Client SonarrAPI
}

// Scrape returns Stats from a Sonarr instance
func (s SonarrScraper) Scrape(ctx context.Context) (Stats, error) {
	stats := Stats{
		URL: s.Client.GetURL(),
	}

	var err error
	stats.Version, err = s.getVersion(ctx)
	if err == nil {
		stats.Health, err = s.getHealth(ctx)
	}
	if err == nil {
		stats.Calendar, err = s.getCalendar(ctx)
	}
	if err == nil {
		stats.Queued, err = s.getQueued(ctx)
	}
	if err == nil {
		stats.Monitored, stats.Unmonitored, err = s.getMonitored(ctx)
	}

	return stats, err
}

func (s SonarrScraper) getVersion(ctx context.Context) (string, error) {
	systemStatus, err := s.Client.GetSystemStatus(ctx)
	if err != nil {
		return "", err
	}
	return systemStatus.Version, nil
}

func (s SonarrScraper) getHealth(ctx context.Context) (map[string]int, error) {
	var healthEntries map[string]int
	health, err := s.Client.GetHealth(ctx)
	if err == nil {
		healthEntries = make(map[string]int)
		for _, entry := range health {
			value := healthEntries[entry.Type]
			healthEntries[entry.Type] = value + 1
		}
	}
	return healthEntries, err
}

func (s SonarrScraper) getCalendar(ctx context.Context) ([]string, error) {
	calendar, err := s.Client.GetCalendar(ctx)
	if err != nil {
		return nil, err
	}
	entries := make([]string, 0, len(calendar))
	for _, entry := range calendar {
		if entry.HasFile {
			continue
		}

		var showName string
		showName, err = s.getShowName(ctx, entry.SeriesID)
		if err != nil {
			return nil, fmt.Errorf("getShowName: %w", err)
		}

		entries = append(entries, fmt.Sprintf("%s - S%02dE%02d - %s",
			showName,
			entry.SeasonNumber,
			entry.EpisodeNumber,
			entry.Title))
	}
	return entries, nil
}

func (s SonarrScraper) getQueued(ctx context.Context) ([]QueuedFile, error) {
	queued, err := s.Client.GetQueue(ctx)
	if err != nil {
		return nil, err
	}
	var entries []QueuedFile
	for _, entry := range queued.Records {
		var episode xxxarr.SonarrEpisodeResponse
		episode, err = s.Client.GetEpisodeByID(ctx, entry.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("GetEpisideByID: %w", err)
		}

		entries = append(entries, QueuedFile{
			Name: fmt.Sprintf("%s - S%02dE%02d - %s",
				episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title),
			TotalBytes:      entry.Size,
			DownloadedBytes: entry.Size - entry.Sizeleft,
		})
	}
	return entries, nil
}

func (s SonarrScraper) getMonitored(ctx context.Context) (int, int, error) {
	var monitored, unmonitored int
	movies, err := s.Client.GetSeries(ctx)
	if err == nil {
		for _, entry := range movies {
			if entry.Monitored {
				monitored++
			} else {
				unmonitored++
			}
		}
	}
	return monitored, unmonitored, err
}

func (s SonarrScraper) getShowName(ctx context.Context, id int) (string, error) {
	var title string
	show, err := s.Client.GetSeriesByID(ctx, id)
	if err == nil {
		title = show.Title
	}
	return title, err
}
