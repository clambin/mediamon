package updater

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
)

// SonarrUpdater collects Stats from a Sonarr instance
type SonarrUpdater struct {
	Client xxxarr.SonarrAPI
}

// GetStats returns Stats from a Sonarr instance
func (su SonarrUpdater) GetStats() (response Stats, err error) {
	var stats Stats

	ctx := context.Background()

	stats.URL = su.Client.GetURL()

	stats.Version, err = su.getVersion(ctx)

	if err == nil {
		stats.Calendar, err = su.getCalendar(ctx)
	}

	if err == nil {
		stats.Queued, err = su.getQueued(ctx)
	}

	if err == nil {
		stats.Monitored, stats.Unmonitored, err = su.getMonitored(ctx)
	}

	return stats, err
}

func (su SonarrUpdater) getVersion(ctx context.Context) (string, error) {
	systemStatus, err := su.Client.GetSystemStatus(ctx)
	if err != nil {
		return "", err
	}
	return systemStatus.Version, nil
}

func (su SonarrUpdater) getCalendar(ctx context.Context) (entries []string, err error) {
	var calendar []xxxarr.SonarrCalendarResponse
	calendar, err = su.Client.GetCalendar(ctx)
	if err != nil {
		return
	}

	for _, entry := range calendar {
		if entry.HasFile {
			continue
		}

		var showName string
		showName, err = su.getShowName(ctx, entry.SeriesID)

		entries = append(entries, fmt.Sprintf("%s - S%02dE%02d - %s",
			showName,
			entry.SeasonNumber,
			entry.EpisodeNumber,
			entry.Title))
	}
	return
}

func (su SonarrUpdater) getQueued(ctx context.Context) (entries []QueuedFile, err error) {
	var queued xxxarr.SonarrQueueResponse
	queued, err = su.Client.GetQueue(ctx)
	if err != nil {
		return
	}

	for _, entry := range queued.Records {
		var episode xxxarr.SonarrEpisodeResponse
		episode, err = su.Client.GetEpisodeByID(ctx, entry.EpisodeID)
		if err != nil {
			return
		}

		entries = append(entries, QueuedFile{
			Name: fmt.Sprintf("%s - S%02dE%02d - %s",
				episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title),
			TotalBytes:      entry.Size,
			DownloadedBytes: entry.Size - entry.Sizeleft,
		})

	}
	return
}

func (su SonarrUpdater) getMonitored(ctx context.Context) (monitored int, unmonitored int, err error) {
	var movies []xxxarr.SonarrSeriesResponse
	movies, err = su.Client.GetSeries(ctx)
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

func (su SonarrUpdater) getShowName(ctx context.Context, id int) (title string, err error) {
	var show xxxarr.SonarrSeriesResponse
	show, err = su.Client.GetSeriesByID(ctx, id)
	if err == nil {
		title = show.Title
	}
	return
}
