package caller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheTable_ShouldCache(t *testing.T) {
	table := CacheTable{Table: []CacheTableEntry{
		{Endpoint: `/foo`},
		{Endpoint: `/foo/[\d+]`, IsRegExp: true},
	}}

	type testcase struct {
		input  string
		expiry time.Duration
		match  bool
	}
	for _, tc := range []testcase{
		{input: "/foo", match: true},
		{input: "/foo/123", match: true},
		{input: "/foo/bar", match: false},
		{input: "/bar", match: false},
	} {
		found, expiry := table.shouldCache(tc.input)
		assert.Equal(t, tc.match, found, tc.input)
		assert.Equal(t, tc.expiry, expiry, tc.input)

	}
}

func TestCacheTable_CacheEverything(t *testing.T) {
	table := CacheTable{}

	found, _ := table.shouldCache("/")
	assert.True(t, found)
}

func TestCacheTable_Invalid_Input(t *testing.T) {
	table := CacheTable{Table: []CacheTableEntry{
		{Endpoint: `/foo/[\d+`, IsRegExp: true},
	}}

	assert.Panics(t, func() { table.shouldCache("/foo") })
}
