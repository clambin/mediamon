package qplex

import (
	"cmp"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestViewCount(t *testing.T) {
	vc := make(viewCount)

	vc.merge(viewCount{
		"1": ViewCountEntry{
			Library: "foo",
			Title:   "foo 1",
			Views:   1,
		},
	})
	vc.merge(viewCount{
		"1": ViewCountEntry{
			Library: "foo",
			Title:   "foo 1",
			Views:   2,
		},
	})
	vc.merge(viewCount{
		"2": ViewCountEntry{
			Library: "bar",
			Title:   "bar 1",
			Views:   1,
		},
	})
	vc.merge(viewCount{
		"3": ViewCountEntry{
			Library: "foo",
			Title:   "foo 2",
			Views:   1,
		},
	})

	f := vc.flatten()
	slices.SortFunc(f, func(a, b ViewCountEntry) int {
		return cmp.Compare(a.Title, b.Title)
	})
	assert.Equal(t, []ViewCountEntry{
		{Library: "bar", Title: "bar 1", Views: 1},
		{Library: "foo", Title: "foo 1", Views: 3},
		{Library: "foo", Title: "foo 2", Views: 1},
	}, f)
}
