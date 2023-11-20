package clients

import (
	"context"
	"fmt"
	xxxarr2 "github.com/clambin/mediaclients/xxxarr"
)

type Sonarr struct {
	Client SonarrClient
}

type SonarrClient interface {
	GetSystemStatus(ctx context.Context) (xxxarr2.SonarrSystemStatus, error)
	GetHealth(ctx context.Context) ([]xxxarr2.SonarrHealth, error)
	GetCalendar(ctx context.Context) ([]xxxarr2.SonarrCalendar, error)
	GetQueue(ctx context.Context) ([]xxxarr2.SonarrQueue, error)
	GetEpisodeByID(ctx context.Context, id int) (xxxarr2.SonarrEpisode, error)
	GetSeries(ctx context.Context) ([]xxxarr2.SonarrSeries, error)
}

func (s Sonarr) GetVersion(ctx context.Context) (string, error) {
	systemStatus, err := s.Client.GetSystemStatus(ctx)
	return systemStatus.Version, err
}

func (s Sonarr) GetHealth(ctx context.Context) (map[string]int, error) {
	health := make(map[string]int)
	healthItems, err := s.Client.GetHealth(ctx)
	for i := range healthItems {
		value := health[healthItems[i].Type]
		health[healthItems[i].Type] = value + 1
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
	entries := make([]QueuedItem, len(queued))
	for i := range queued {
		var episode xxxarr2.SonarrEpisode
		episode, err = s.Client.GetEpisodeByID(ctx, queued[i].EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("GetEpisideByID: %w", err)
		}

		entries[i] = QueuedItem{
			Name: fmt.Sprintf("%s - S%02dE%02d - %s",
				episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title),
			TotalBytes:      queued[i].Size,
			DownloadedBytes: queued[i].Size - queued[i].SizeLeft,
		}
	}
	return entries, err
}

func (s Sonarr) GetLibrary(ctx context.Context) (Library, error) {
	var library Library
	series, err := s.Client.GetSeries(ctx)
	for i := range series {
		if series[i].Monitored {
			library.Monitored++
		} else {
			library.Unmonitored++
		}
	}
	return library, err
}
