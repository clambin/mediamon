package clients

import (
	"context"
	xxxarr2 "github.com/clambin/mediaclients/xxxarr"
)

type RadarrClient interface {
	GetSystemStatus(ctx context.Context) (xxxarr2.RadarrSystemStatusResponse, error)
	GetHealth(ctx context.Context) ([]xxxarr2.RadarrHealthResponse, error)
	GetCalendar(ctx context.Context) ([]xxxarr2.RadarrCalendarResponse, error)
	GetQueue(ctx context.Context) (xxxarr2.RadarrQueueResponse, error)
	GetMovies(ctx context.Context) ([]xxxarr2.RadarrMovieResponse, error)
}

//var _ xxxarr.XXXArrGetter = &Radarr{}

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
	for _, item := range healthItems {
		value := health[item.Type]
		health[item.Type] = value + 1
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
	var entries []QueuedItem
	for _, entry := range queued.Records {
		entries = append(entries, QueuedItem{
			Name:            entry.Title,
			TotalBytes:      int64(entry.Size),
			DownloadedBytes: int64(entry.Size) - int64(entry.Sizeleft),
		})
	}
	return entries, err
}

func (r Radarr) GetLibrary(ctx context.Context) (Library, error) {
	var library Library
	series, err := r.Client.GetMovies(ctx)
	for _, entry := range series {
		if entry.Monitored {
			library.Monitored++
		} else {
			library.Unmonitored++
		}
	}
	return library, err
}
