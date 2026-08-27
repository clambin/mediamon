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
	movieCountMetric = prometheus.NewDesc(
		"mediamon_plex_movie_count",
		"Total number of movies in Plex library",
		[]string{"url", "library"},
		nil,
	)
	showCountMetric = prometheus.NewDesc(
		"mediamon_plex_show_count",
		"Total number of shows in Plex library",
		[]string{"url", "library"},
		nil,
	)
	episodeCountMetric = prometheus.NewDesc(
		"mediamon_plex_episode_count",
		"Total number of episodes in Plex library",
		[]string{"url", "library"},
		nil,
	)
)

type libraryGetter interface {
	GetLibraries(ctx context.Context) ([]plex.Library, error)
	GetAllLibraryMedia(ctx context.Context, key string) ([]plex.MediaMetadata, error)
}

type libraryCollector struct {
	libraryGetter
	logger *slog.Logger
	url    string
	measurer.CachingMeasurer[map[string]libraryInfo]
}

func newLibraryCollector(client libraryGetter, url string, logger *slog.Logger) prometheus.Collector {
	c := &libraryCollector{
		libraryGetter: client,
		url:           url,
		logger:        logger,
	}
	c.CachingMeasurer = measurer.CachingMeasurer[map[string]libraryInfo]{
		Do:       c.analyzeLibraries,
		Interval: libraryRefreshInterval,
	}
	return c
}

func (c *libraryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- libraryBytesMetric
	ch <- movieCountMetric
	ch <- showCountMetric
	ch <- episodeCountMetric
}

func (c *libraryCollector) Collect(ch chan<- prometheus.Metric) {
	libraryInfos, err := c.Measure(context.Background())
	if err != nil {
		c.logger.Error("fail to collect library metrics", "err", err)
		return
	}

	for name, info := range libraryInfos {
		if info.movies > 0 {
			// report movie count
			ch <- prometheus.MustNewConstMetric(movieCountMetric, prometheus.GaugeValue, float64(info.movies), c.url, name)
		}
		if info.shows > 0 {
			// report show count
			ch <- prometheus.MustNewConstMetric(showCountMetric, prometheus.GaugeValue, float64(info.shows), c.url, name)
			// report episode count
			ch <- prometheus.MustNewConstMetric(episodeCountMetric, prometheus.GaugeValue, float64(info.episodes), c.url, name)
		}
		// report total disk space
		ch <- prometheus.MustNewConstMetric(libraryBytesMetric, prometheus.GaugeValue, float64(info.totalSize), c.url, name)
	}
}

type libraryInfo struct {
	movies    int
	shows     int
	episodes  int
	totalSize int64
}

func (c *libraryCollector) analyzeLibraries(ctx context.Context) (map[string]libraryInfo, error) {
	libraries, err := c.GetLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetLibraries: %w", err)
	}

	result := make(map[string]libraryInfo, len(libraries))
	for index := range libraries {
		media, err := c.GetAllLibraryMedia(ctx, libraries[index].Key)
		if err != nil {
			return nil, fmt.Errorf("GetAllLibraryMedia (%s): %w", libraries[index].Type, err)
		}
		var info libraryInfo
		titles := make(map[string]struct{})
		for _, entry := range media {
			switch libraries[index].Type {
			case "movie":
				titles[entry.Title] = struct{}{}
			case "show":
				titles[entry.GrandparentTitle] = struct{}{}
				info.episodes++
			}
			info.totalSize += getMediaSize(entry.Media)
		}
		switch libraries[index].Type {
		case "show":
			info.shows = len(titles)
		case "movie":
			info.movies = len(titles)
		}
		result[libraries[index].Title] = info
	}
	return result, nil
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
