package clients

import (
	"context"
	xxxarr2 "github.com/clambin/mediaclients/xxxarr"
)

type RadarrClient interface {
	GetSystemStatus(ctx context.Context) (xxxarr2.RadarrSystemStatus, error)
	GetHealth(ctx context.Context) ([]xxxarr2.RadarrHealth, error)
	GetCalendar(ctx context.Context) ([]xxxarr2.RadarrCalendar, error)
	GetQueue(ctx context.Context) ([]xxxarr2.RadarrQueue, error)
	GetMovies(ctx context.Context) ([]xxxarr2.RadarrMovie, error)
}

type Radarr struct {
	Client RadarrClient
}

func (r Radarr) GetVersion(ctx context.Context) (string, error) {
	systemStatus, err := r.Client.GetSystemStatus(ctx)
	return systemStatus.Version, err
}

func (r Radarr) GetHealth(ctx context.Context) (map[string]int, error) {
	health := make(map[string]int)
	healthItems, err := r.Client.GetHealth(ctx)
	for i := range healthItems {
		healthType := healthItems[i].Type
		value := health[healthType]
		health[healthType] = value + 1
	}
	return health, err
}

func (r Radarr) GetCalendar(ctx context.Context) ([]string, error) {
	items, err := r.Client.GetCalendar(ctx)
	calendar := make([]string, len(items))
	for i := range items {
		calendar[i] = items[i].Title
	}
	return calendar, err
}

func (r Radarr) GetQueue(ctx context.Context) ([]QueuedItem, error) {
	queued, err := r.Client.GetQueue(ctx)
	entries := make([]QueuedItem, len(queued))
	for i := range queued {
		entries[i] = QueuedItem{
			Name:            queued[i].Title,
			TotalBytes:      queued[i].Size,
			DownloadedBytes: queued[i].Size - queued[i].SizeLeft,
		}
	}
	return entries, err
}

func (r Radarr) GetLibrary(ctx context.Context) (Library, error) {
	var library Library
	series, err := r.Client.GetMovies(ctx)
	for i := range series {
		if series[i].Monitored {
			library.Monitored++
		} else {
			library.Unmonitored++
		}
	}
	return library, err
}
