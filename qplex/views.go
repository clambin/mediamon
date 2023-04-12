package qplex

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"golang.org/x/exp/slog"
	"sort"
	"time"
)

type PlexGetter interface {
	SetAuthToken(string)
	GetLibraries(context.Context) (plex.Libraries, error)
	GetMovieLibrary(context.Context, string) (plex.MovieLibrary, error)
	GetShowLibrary(context.Context, string) (plex.ShowLibrary, error)
}

func GetViews(ctx context.Context, c PlexGetter, tokens []string, reverse bool) ([]ViewCountEntry, error) {
	totalViewCount, err := getViewCount(ctx, c, tokens)
	if err != nil {
		return nil, err
	}

	flattened := totalViewCount.flatten()
	sort.Slice(flattened, func(i, j int) bool {
		if reverse {
			return flattened[i].Views > flattened[j].Views
		}
		return flattened[i].Views < flattened[j].Views
	})
	return flattened, nil
}

func getViewCount(ctx context.Context, c PlexGetter, authTokens []string) (viewCount, error) {
	totalViewCount := make(viewCount)

	for _, token := range authTokens {
		c.SetAuthToken(token)
		vc, err := getUserViewCount(ctx, c)
		if err != nil {
			return nil, err
		}
		totalViewCount.merge(vc)
	}
	return totalViewCount, nil
}

func getUserViewCount(ctx context.Context, c PlexGetter) (viewCount, error) {
	result := make(viewCount)
	libraries, err := c.GetLibraries(ctx)
	if err != nil {
		return nil, err
	}
	for _, library := range libraries.Directory {
		var libraryViewCount viewCount
		switch library.Type {
		case "movie":
			libraryViewCount, err = getMovieViewCount(ctx, c, library)
		case "show":
			libraryViewCount, err = getShowViewCount(ctx, c, library)
		}
		if err != nil {
			return nil, err
		}

		result.merge(libraryViewCount)
	}
	return result, err
}

func getMovieViewCount(ctx context.Context, client PlexGetter, library plex.LibrariesDirectory) (viewCount, error) {
	result := make(viewCount)
	entries, err := client.GetMovieLibrary(ctx, library.Key)
	if err != nil {
		return result, err
	}
	for _, entry := range entries.Metadata {
		existing, ok := result[guid(entry.Guid)]
		if !ok {
			existing = ViewCountEntry{
				Library: library.Title,
				Title:   entry.Title,
			}
		}
		existing.Views += entry.ViewCount
		result[guid(entry.Guid)] = existing
	}
	return result, nil
}

func getShowViewCount(ctx context.Context, client PlexGetter, library plex.LibrariesDirectory) (viewCount, error) {
	result := make(viewCount)
	entries, err := client.GetShowLibrary(ctx, library.Key)
	if err != nil {
		return result, err
	}
	for _, entry := range entries.Metadata {
		existing, ok := result[guid(entry.Guid)]
		if !ok {
			existing = ViewCountEntry{
				Library: library.Title,
				Title:   entry.Title,
			}
		}
		existing.Views += entry.ViewCount
		result[guid(entry.Guid)] = existing

		slog.Debug("show found",
			slog.String("name", entry.Title),
			slog.Time("added", time.Time(entry.AddedAt)),
			slog.Time("updated", time.Time(entry.UpdatedAt)),
			slog.Time("viewed", time.Time(entry.LastViewedAt)),
		)
	}
	return result, nil
}
