package scraper

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
)

// SonarrScraper collects Stats from a Sonarr instance
type SonarrScraper struct {
	Client xxxarr.SonarrAPI
}

var _ Scraper = &SonarrScraper{}

// Scrape returns Stats from a Sonarr instance
func (s SonarrScraper) Scrape() (response Stats, err error) {
	var stats Stats

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

	return stats, err
}

func (s SonarrScraper) getVersion(ctx context.Context) (version string, err error) {
	var systemStatus xxxarr.SonarrSystemStatusResponse
	systemStatus, err = s.Client.GetSystemStatus(ctx)
	if err == nil {
		version = systemStatus.Version
	}
	return
}

func (s SonarrScraper) getCalendar(ctx context.Context) (entries []string, err error) {
	var calendar []xxxarr.SonarrCalendarResponse
	calendar, err = s.Client.GetCalendar(ctx)
	if err == nil {
		for _, entry := range calendar {
			if entry.HasFile {
				continue
			}

			var showName string
			showName, err = s.getShowName(ctx, entry.SeriesID)

			entries = append(entries, fmt.Sprintf("%s - S%02dE%02d - %s",
				showName,
				entry.SeasonNumber,
				entry.EpisodeNumber,
				entry.Title))
		}
	}
	return
}

func (s SonarrScraper) getQueued(ctx context.Context) (entries []QueuedFile, err error) {
	var queued xxxarr.SonarrQueueResponse
	queued, err = s.Client.GetQueue(ctx)
	if err == nil {
		for _, entry := range queued.Records {
			var episode xxxarr.SonarrEpisodeResponse
			episode, err = s.Client.GetEpisodeByID(ctx, entry.EpisodeID)
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
	}
	return
}

func (s SonarrScraper) getMonitored(ctx context.Context) (monitored int, unmonitored int, err error) {
	var movies []xxxarr.SonarrSeriesResponse
	movies, err = s.Client.GetSeries(ctx)
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

func (s SonarrScraper) getShowName(ctx context.Context, id int) (title string, err error) {
	var show xxxarr.SonarrSeriesResponse
	show, err = s.Client.GetSeriesByID(ctx, id)
	if err == nil {
		title = show.Title
	}
	return
}
