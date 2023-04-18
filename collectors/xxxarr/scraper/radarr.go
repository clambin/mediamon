package scraper

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/xxxarr"
)

// RadarrScraper collects Stats from a Radarr instance
type RadarrScraper struct {
	Client xxxarr.RadarrAPI
}

// Scrape returns Stats from a Radarr instance
func (s RadarrScraper) Scrape(ctx context.Context) (Stats, error) {
	stats := Stats{
		URL: s.Client.GetURL(),
	}

	var err error
	stats.URL = s.Client.GetURL()

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

func (s RadarrScraper) getVersion(ctx context.Context) (string, error) {
	systemStatus, err := s.Client.GetSystemStatus(ctx)
	return systemStatus.Version, err
}

func (s RadarrScraper) getHealth(ctx context.Context) (map[string]int, error) {
	health, err := s.Client.GetHealth(ctx)
	if err != nil {
		return nil, err
	}
	healthEntries := make(map[string]int)
	for _, entry := range health {
		value := healthEntries[entry.Type]
		healthEntries[entry.Type] = value + 1
	}
	return healthEntries, nil
}

func (s RadarrScraper) getCalendar(ctx context.Context) ([]string, error) {
	calendar, err := s.Client.GetCalendar(ctx)
	if err != nil {
		return nil, err
	}
	var movieTitles []string
	for _, entry := range calendar {
		if !entry.HasFile {
			movieTitles = append(movieTitles, entry.Title)
		}
	}
	return movieTitles, nil
}

func (s RadarrScraper) getQueued(ctx context.Context) ([]QueuedFile, error) {
	queued, err := s.Client.GetQueue(ctx)
	if err != nil {
		return nil, err
	}
	var queuedFiles []QueuedFile
	for _, entry := range queued.Records {
		var movie xxxarr.RadarrMovieResponse
		movie, err = s.Client.GetMovieByID(ctx, entry.MovieID)
		if err != nil {
			return nil, err
		}

		queuedFiles = append(queuedFiles, QueuedFile{
			Name:            movie.Title,
			TotalBytes:      float64(entry.Size),
			DownloadedBytes: float64(entry.Size - entry.Sizeleft),
		})
	}
	return queuedFiles, nil
}

func (s RadarrScraper) getMonitored(ctx context.Context) (int, int, error) {
	movies, err := s.Client.GetMovies(ctx)
	if err != nil {
		return 0, 0, err
	}
	var monitored, unmonitored int
	for _, entry := range movies {
		if entry.Monitored {
			monitored++
		} else {
			unmonitored++
		}
	}
	return monitored, unmonitored, nil
}
