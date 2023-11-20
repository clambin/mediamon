package clients

import (
	"context"
	"fmt"
	xxxarr2 "github.com/clambin/mediaclients/xxxarr"
)

//var _ xxxarr.XXXArrGetter = &Sonarr{}

type Sonarr struct {
	Client SonarrClient
}

type SonarrClient interface {
	GetSystemStatus(ctx context.Context) (xxxarr2.SonarrSystemStatusResponse, error)
	GetHealth(ctx context.Context) ([]xxxarr2.SonarrHealthResponse, error)
	GetCalendar(ctx context.Context) ([]xxxarr2.SonarrCalendarResponse, error)
	GetQueue(ctx context.Context) (xxxarr2.SonarrQueueResponse, error)
	GetEpisodeByID(ctx context.Context, id int) (xxxarr2.SonarrEpisodeResponse, error)
	GetSeries(ctx context.Context) ([]xxxarr2.SonarrSeriesResponse, error)
}

func (s Sonarr) GetVersion(ctx context.Context) (string, error) {
	systemStatus, err := s.Client.GetSystemStatus(ctx)
	return systemStatus.Version, err
}

func (s Sonarr) GetHealth(ctx context.Context) (map[string]int, error) {
	health := make(map[string]int)
	healthItems, err := s.Client.GetHealth(ctx)
	for _, item := range healthItems {
		value := health[item.Type]
		health[item.Type] = value + 1
	}
	return health, err
}

func (s Sonarr) GetCalendar(ctx context.Context) ([]string, error) {
	items, err := s.Client.GetCalendar(ctx)
	calendar := make([]string, len(items))
	for i := range items {
		calendar[i] = items[i].Title
	}
	return calendar, err
}

func (s Sonarr) GetQueue(ctx context.Context) ([]QueuedItem, error) {
	queued, err := s.Client.GetQueue(ctx)
	var entries []QueuedItem
	for _, entry := range queued.Records {
		var episode xxxarr2.SonarrEpisodeResponse
		episode, err = s.Client.GetEpisodeByID(ctx, entry.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("GetEpisideByID: %w", err)
		}

		entries = append(entries, QueuedItem{
			Name: fmt.Sprintf("%s - S%02dE%02d - %s",
				episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title),
			TotalBytes:      int64(entry.Size),
			DownloadedBytes: int64(entry.Size) - int64(entry.Sizeleft),
		})
	}
	return entries, err
}

func (s Sonarr) GetLibrary(ctx context.Context) (Library, error) {
	var library Library
	series, err := s.Client.GetSeries(ctx)
	for _, entry := range series {
		if entry.Monitored {
			library.Monitored++
		} else {
			library.Unmonitored++
		}
	}
	return library, err
}
