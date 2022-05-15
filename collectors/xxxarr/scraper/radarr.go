package scraper

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
)

// RadarrScraper collects Stats from a Radarr instance
type RadarrScraper struct {
	Client xxxarr.RadarrAPI
}

var _ Scraper = &RadarrScraper{}

// Scrape returns Stats from a Radarr instance
func (s RadarrScraper) Scrape() (stats Stats, err error) {
	ctx := context.Background()

	stats.URL = s.Client.GetURL()

	stats.Version, err = s.getVersion(ctx)
	if err == nil {
		stats.Calendar, err = s.getCalendar(ctx)
	}

	if err == nil {
		stats.Queued, err = s.getQueued(ctx)
	}

	if err == nil {
		stats.Monitored, stats.Unmonitored, err = s.getMonitored(ctx)
	}

	return
}

func (s RadarrScraper) getVersion(ctx context.Context) (string, error) {
	systemStatus, err := s.Client.GetSystemStatus(ctx)
	return systemStatus.Version, err
}

func (s RadarrScraper) getCalendar(ctx context.Context) (entries []string, err error) {
	var calendar []xxxarr.RadarrCalendarResponse
	calendar, err = s.Client.GetCalendar(ctx)
	if err == nil {
		for _, entry := range calendar {
			if !entry.HasFile {
				entries = append(entries, entry.Title)
			}
		}
	}
	return
}

func (s RadarrScraper) getQueued(ctx context.Context) (entries []QueuedFile, err error) {
	var queued xxxarr.RadarrQueueResponse
	queued, err = s.Client.GetQueue(ctx)
	if err == nil {
		for _, entry := range queued.Records {
			var movie xxxarr.RadarrMovieResponse
			movie, err = s.Client.GetMovieByID(ctx, entry.MovieID)
			if err != nil {
				return
			}

			entries = append(entries, QueuedFile{
				Name:            movie.Title,
				TotalBytes:      float64(entry.Size),
				DownloadedBytes: float64(entry.Size - entry.Sizeleft),
			})
		}
	}
	return
}

func (s RadarrScraper) getMonitored(ctx context.Context) (monitored int, unmonitored int, err error) {
	var movies []xxxarr.RadarrMovieResponse
	movies, err = s.Client.GetMovies(ctx)
	if err == nil {
		for _, entry := range movies {
			if entry.Monitored {
				monitored++
			} else {
				unmonitored++
			}
		}
	}
	return
}
