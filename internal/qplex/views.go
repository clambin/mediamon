package qplex

import (
	"cmp"
	"context"
	"github.com/clambin/mediaclients/plex"
	"slices"
)

type PlexGetter interface {
	SetAuthToken(string)
	GetLibraries(context.Context) ([]plex.Library, error)
	GetMovies(context.Context, string) ([]plex.Movie, error)
	GetShows(context.Context, string) ([]plex.Show, error)
}

func GetViews(ctx context.Context, c PlexGetter, tokens []string, reverse bool) ([]ViewCountEntry, error) {
	totalViewCount, err := getViewCount(ctx, c, tokens)
	if err != nil {
		return nil, err
	}

	flattened := totalViewCount.flatten()
	slices.SortFunc(flattened, func(a, b ViewCountEntry) int {
		if reverse {
			a, b = b, a
		}
		return cmp.Compare(a.Views, b.Views)
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
	for _, library := range libraries {
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

func getMovieViewCount(ctx context.Context, client PlexGetter, library plex.Library) (viewCount, error) {
	result := make(viewCount)
	entries, err := client.GetMovies(ctx, library.Key)
	if err != nil {
		return result, err
	}
	for _, entry := range entries {
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

func getShowViewCount(ctx context.Context, client PlexGetter, library plex.Library) (viewCount, error) {
	result := make(viewCount)
	entries, err := client.GetShows(ctx, library.Key)
	if err != nil {
		return result, err
	}
	for _, entry := range entries {
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
