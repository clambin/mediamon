package plex

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/plex"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

var libraryMetric = prometheus.NewDesc(
	prometheus.BuildFQName("mediamon", "plex", "library_entry_bytes"),
	"Library file sizes",
	[]string{"url", "library", "title"},
	nil,
)

type libraryCollector struct {
	libraryGetter
	url string
	l   *slog.Logger
}

type libraryGetter interface {
	GetLibraries(ctx context.Context) ([]plex.Library, error)
	GetMovies(ctx context.Context, key string) ([]plex.Movie, error)
	GetShows(ctx context.Context, key string) ([]plex.Show, error)
	GetSeasons(ctx context.Context, key string) ([]plex.Season, error)
	GetEpisodes(ctx context.Context, key string) ([]plex.Episode, error)
}

func (c libraryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- libraryMetric
}

func (c libraryCollector) Collect(ch chan<- prometheus.Metric) {
	sizes, err := c.reportSizes()
	if err != nil {
		c.l.Error("failed to collect plex library stats", "err", err)
		return
	}

	for library, entries := range sizes {
		for _, entry := range entries {
			ch <- prometheus.MustNewConstMetric(libraryMetric, prometheus.GaugeValue, float64(entry.size), c.url, library, entry.title)
		}
	}
}

type libraryEntry struct {
	title string
	size  int64
}

func (c libraryCollector) reportSizes() (map[string][]libraryEntry, error) {
	ctx := context.Background()
	libraries, err := c.libraryGetter.GetLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetLibraries: %w", err)
	}

	result := make(map[string][]libraryEntry)
	var sizes []libraryEntry
	for index := range libraries {
		switch libraries[index].Type {
		case "movie":
			if sizes, err = c.getMovieTotals(ctx, libraries[index].Key); err != nil {
				return nil, fmt.Errorf("getMovieTotals: %w", err)
			}
		case "show":
			if sizes, err = c.getShowTotals(ctx, libraries[index].Key); err != nil {
				return nil, fmt.Errorf("getShowTotals: %w", err)
			}
		}
		result[libraries[index].Title] = sizes
	}
	return result, nil
}

func (c libraryCollector) getMovieTotals(ctx context.Context, key string) ([]libraryEntry, error) {
	movies, err := c.GetMovies(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("GetMovies: %w", err)
	}
	entries := make([]libraryEntry, 0, len(movies))
	for index := range movies {
		entries = append(entries, libraryEntry{
			title: movies[index].Title,
			size:  getMediaSize(movies[index].Media),
		})
	}
	return entries, nil
}

func (c libraryCollector) getShowTotals(ctx context.Context, key string) ([]libraryEntry, error) {
	shows, err := c.GetShows(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("GetShows: %w", err)
	}

	entries := make([]libraryEntry, 0, len(shows))
	for index := range shows {
		size, err := c.getShowTotal(ctx, shows[index].RatingKey)
		if err != nil {
			return nil, fmt.Errorf("getShowTotal: %w", err)
		}
		if size > 0 {
			entries = append(entries, libraryEntry{
				title: shows[index].Title,
				size:  size,
			})
		}
	}
	return entries, nil
}

func (c libraryCollector) getShowTotal(ctx context.Context, key string) (int64, error) {
	seasons, err := c.GetSeasons(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("GetSeasons: %w", err)
	}
	var size int64
	for index := range seasons {
		episodes, err := c.GetEpisodes(ctx, seasons[index].RatingKey)
		if err != nil {
			return 0, fmt.Errorf("GetEpisodes: %w", err)
		}
		for index2 := range episodes {
			size += getMediaSize(episodes[index2].Media)
		}
	}
	return size, nil
}

func getMediaSize(medias []plex.Media) int64 {
	for _, media := range medias {
		for _, part := range media.Part {
			if part.Size > 0 {
				return part.Size
			}
		}
	}
	return 0
}
