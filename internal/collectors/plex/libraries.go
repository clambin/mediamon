package plex

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/measurer"
	"github.com/prometheus/client_golang/prometheus"
)

const libraryRefreshInterval = 6 * time.Hour

var (
	libraryBytesMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "library_bytes"),
		"Library size in bytes",
		[]string{"url", "library"},
		nil,
	)
	libraryCountMetric = prometheus.NewDesc(
		prometheus.BuildFQName("mediamon", "plex", "library_count"),
		"Library size in number of entries",
		[]string{"url", "library"},
		nil,
	)
)

type libraryGetter interface {
	GetLibraries(ctx context.Context) ([]plex.Library, error)
	GetMovies(ctx context.Context, key string) ([]plex.Movie, error)
	GetShows(ctx context.Context, key string) ([]plex.Show, error)
	GetSeasons(ctx context.Context, key string) ([]plex.Season, error)
	GetEpisodes(ctx context.Context, key string) ([]plex.Episode, error)
}

type libraryCollector struct {
	libraryGetter
	logger *slog.Logger
	url    string
	measurer.Cached[map[string][]libraryEntry]
}

func newLibraryCollector(client libraryGetter, url string, logger *slog.Logger) prometheus.Collector {
	c := &libraryCollector{
		libraryGetter: client,
		url:           url,
		logger:        logger,
	}
	c.Cached = measurer.Cached[map[string][]libraryEntry]{
		Do:       c.getLibraries,
		Interval: libraryRefreshInterval,
	}
	return c
}

func (c *libraryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- libraryBytesMetric
	ch <- libraryCountMetric
}

func (c *libraryCollector) Collect(ch chan<- prometheus.Metric) {
	libraries, err := c.Measure(context.Background())
	if err != nil {
		c.logger.Error("fail to collect library metrics", "err", err)
		return
	}

	for library, entries := range libraries {
		ch <- prometheus.MustNewConstMetric(libraryCountMetric, prometheus.GaugeValue, float64(len(entries)), c.url, library)
		var size int64
		for _, entry := range entries {
			size += entry.size
		}
		ch <- prometheus.MustNewConstMetric(libraryBytesMetric, prometheus.GaugeValue, float64(size), c.url, library)
	}
}

type libraryEntry struct {
	title string
	size  int64
}

func (c *libraryCollector) getLibraries(ctx context.Context) (map[string][]libraryEntry, error) {
	libraries, err := c.GetLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetLibraries: %w", err)
	}

	result := make(map[string][]libraryEntry)
	var entries []libraryEntry
	for index := range libraries {
		switch libraries[index].Type {
		case "movie":
			entries, err = c.getMovieTotals(ctx, libraries[index].Key)
		case "show":
			entries, err = c.getShowTotals(ctx, libraries[index].Key)
		}
		if err != nil {
			return nil, fmt.Errorf("getTotals (%s): %w", libraries[index].Type, err)
		}
		result[libraries[index].Title] = entries
	}
	return result, nil
}

func (c *libraryCollector) getMovieTotals(ctx context.Context, key string) ([]libraryEntry, error) {
	movies, err := c.GetMovies(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("GetMovies: %w", err)
	}
	entries := make([]libraryEntry, len(movies))
	for index := range movies {
		entries[index] = libraryEntry{
			title: movies[index].Title,
			size:  getMediaSize(movies[index].Media),
		}
	}
	return entries, nil
}

func (c *libraryCollector) getShowTotals(ctx context.Context, key string) ([]libraryEntry, error) {
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

func (c *libraryCollector) getShowTotal(ctx context.Context, key string) (int64, error) {
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
