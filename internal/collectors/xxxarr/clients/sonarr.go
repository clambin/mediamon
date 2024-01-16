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
	calendar, err := s.Client.GetCalendar(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetCalendar: %w", err)
	}
	ids := getIDs(calendar, func(queue xxxarr2.SonarrCalendar) int { return queue.ID })
	entries := make([]string, len(calendar))
	for i := range calendar {
		var name string
		name, err = s.getEpisodeName(ctx, ids[i])
		if err != nil {
			return nil, err
		}
		entries[i] = name
	}
	return entries, err
}

func (s Sonarr) GetQueue(ctx context.Context) ([]QueuedItem, error) {
	queued, err := s.Client.GetQueue(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetQueue: %w", err)
	}
	ids := getIDs(queued, func(queue xxxarr2.SonarrQueue) int { return queue.EpisodeID })
	entries := make([]QueuedItem, len(queued))
	for i := range queued {
		var name string
		name, err = s.getEpisodeName(ctx, ids[i])
		if err != nil {
			return nil, err
		}
		entries[i] = QueuedItem{
			Name:            name,
			TotalBytes:      queued[i].Size,
			DownloadedBytes: queued[i].Size - queued[i].SizeLeft,
		}
	}
	return entries, err
}

func (s Sonarr) getEpisodeName(ctx context.Context, id int) (string, error) {
	episode, err := s.Client.GetEpisodeByID(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s - S%02dE%02d - %s", episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title), nil
}

func getIDs[T ~[]E, E any](episodes T, getID func(E) int) []int {
	ids := make([]int, len(episodes))
	for i := range episodes {
		ids[i] = getID(episodes[i])
	}
	return ids
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
