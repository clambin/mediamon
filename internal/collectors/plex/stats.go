package plex

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"time"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/measurer"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	movieCountMetric = prometheus.NewDesc(
		"mediamon_plex_movie_count",
		"Total number of movies in Plex library",
		[]string{"url"},
		nil,
	)
	showCountMetric = prometheus.NewDesc(
		"mediamon_plex_show_count",
		"Total number of shows in Plex library",
		[]string{"url"},
		nil,
	)
	episodeCountMetric = prometheus.NewDesc(
		"mediamon_plex_episode_count",
		"Total number of episodes in Plex library",
		[]string{"url"},
		nil,
	)
)

type statsGetter interface {
	GetLibraries(ctx context.Context) ([]plex.Library, error)
	GetMovies(ctx context.Context, key string) ([]plex.Movie, error)
	GetShows(ctx context.Context, key string) ([]plex.Show, error)
	GetSeasons(ctx context.Context, key string) ([]plex.Season, error)
	GetEpisodes(ctx context.Context, key string) ([]plex.Episode, error)
}

var _ prometheus.Collector = &statsCollector{}

type statsCollector struct {
	client     statsGetter
	movieStats measurer.CachingMeasurer[int]
	showStats  measurer.CachingMeasurer[[]int]
	url        string
}

func newStatsCollector(client statsGetter, url string, _ *slog.Logger) *statsCollector {
	const statsCacheInterval = time.Hour
	c := statsCollector{client: client, url: url}
	c.movieStats = measurer.CachingMeasurer[int]{
		Interval: statsCacheInterval,
		Do:       c.getMovieStats,
	}
	c.showStats = measurer.CachingMeasurer[[]int]{
		Interval: statsCacheInterval,
		Do:       c.getShowStats,
	}
	return &c
}

func (s *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- movieCountMetric
	ch <- showCountMetric
	ch <- episodeCountMetric
}

func (s *statsCollector) Collect(ch chan<- prometheus.Metric) {
	movieCount, _ := s.movieStats.Measure(context.Background())
	ch <- prometheus.MustNewConstMetric(movieCountMetric, prometheus.GaugeValue, float64(movieCount), s.url)
	showStats, _ := s.showStats.Measure(context.Background())
	ch <- prometheus.MustNewConstMetric(showCountMetric, prometheus.GaugeValue, float64(len(showStats)), s.url)
	var episodes int
	for _, showStat := range showStats {
		episodes += showStat
	}
	ch <- prometheus.MustNewConstMetric(episodeCountMetric, prometheus.GaugeValue, float64(episodes), s.url)
}

func (s *statsCollector) getMovieStats(ctx context.Context) (int, error) {
	libraries, err := s.client.GetLibraries(ctx)
	if err != nil {
		return 0, fmt.Errorf("get libraries: %w", err)
	}
	var movieCount int
	for _, library := range libraries {
		if library.Type != "movie" {
			continue
		}
		movies, err := s.client.GetMovies(ctx, library.Key)
		if err != nil {
			return 0, fmt.Errorf("get movies %q: %w", library.Key, err)
		}
		movieCount += len(movies)
	}
	return movieCount, nil
}

func (s *statsCollector) getShowStats(ctx context.Context) ([]int, error) {
	libraries, err := s.client.GetLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("get libraries: %w", err)
	}
	showStats := make(map[string]int)
	for _, library := range libraries {
		if library.Type != "show" {
			continue
		}
		shows, err := s.client.GetShows(ctx, library.Key)
		if err != nil {
			return nil, fmt.Errorf("get shows")
		}
		for _, show := range shows {
			var episodeCount int
			seasons, err := s.client.GetSeasons(ctx, show.RatingKey)
			if err != nil {
				return nil, fmt.Errorf("get seasons %v: %w", show.Title, err)
			}
			for _, season := range seasons {
				episodes, err := s.client.GetEpisodes(ctx, season.RatingKey)
				if err != nil {
					return nil, fmt.Errorf("get episodes %v / %v: %w", show.Title, season.Title, err)
				}
				episodeCount += len(episodes)
			}
			showStats[show.Title] += episodeCount
		}
	}
	return slices.Collect(maps.Values(showStats)), nil
}
