package qplex

import (
	"github.com/stretchr/testify/assert"
	"sort"
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
	sort.Slice(f, func(i, j int) bool {
		return f[i].Title < f[j].Title
	})
	assert.Equal(t, []ViewCountEntry{
		{Library: "bar", Title: "bar 1", Views: 1},
		{Library: "foo", Title: "foo 1", Views: 3},
		{Library: "foo", Title: "foo 2", Views: 1},
	}, f)
}
