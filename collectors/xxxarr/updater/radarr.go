package updater

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
)

// RadarrUpdater collects Stats from a Radarr instance
type RadarrUpdater struct {
	Client xxxarr.RadarrAPI
}

var _ StatsGetter = &RadarrUpdater{}

// GetStats returns Stats from a Radarr instance
func (ru RadarrUpdater) GetStats() (stats Stats, err error) {
	ctx := context.Background()

	stats.URL = ru.Client.GetURL()

	stats.Version, err = ru.getVersion(ctx)
	if err == nil {
		stats.Calendar, err = ru.getCalendar(ctx)
	}

	if err == nil {
		stats.Queued, err = ru.getQueued(ctx)
	}

	if err == nil {
		stats.Monitored, stats.Unmonitored, err = ru.getMonitored(ctx)
	}

	return
}

func (ru RadarrUpdater) getVersion(ctx context.Context) (string, error) {
	systemStatus, err := ru.Client.GetSystemStatus(ctx)
	return systemStatus.Version, err
}

func (ru RadarrUpdater) getCalendar(ctx context.Context) (entries []string, err error) {
	var calendar []xxxarr.RadarrCalendarResponse
	calendar, err = ru.Client.GetCalendar(ctx)
	if err != nil {
		return
	}

	for _, entry := range calendar {
		if !entry.HasFile {
			entries = append(entries, entry.Title)
		}
	}
	return
}

func (ru RadarrUpdater) getQueued(ctx context.Context) (entries []QueuedFile, err error) {
	var queued xxxarr.RadarrQueueResponse
	queued, err = ru.Client.GetQueue(ctx)
	if err != nil {
		return
	}

	for _, entry := range queued.Records {
		var movie xxxarr.RadarrMovieResponse
		movie, err = ru.Client.GetMovieByID(ctx, entry.MovieID)
		if err != nil {
			return
		}

		entries = append(entries, QueuedFile{
			Name:            movie.Title,
			TotalBytes:      float64(entry.Size),
			DownloadedBytes: float64(entry.Size - entry.Sizeleft),
		})
	}
	return
}

func (ru RadarrUpdater) getMonitored(ctx context.Context) (monitored int, unmonitored int, err error) {
	var movies []xxxarr.RadarrMovieResponse
	movies, err = ru.Client.GetMovies(ctx)
	if err != nil {
		return
	}

	for _, entry := range movies {
		if entry.Monitored {
			monitored++
		} else {
			unmonitored++
		}
	}

	return
}
