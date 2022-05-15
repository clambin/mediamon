package caller

import (
	"fmt"
	"regexp"
	"time"
)

// CacheTable holds the Endpoints that should be cached
type CacheTable struct {
	Table    []CacheTableEntry
	compiled bool
}

// CacheTableEntry contains a single endpoint that should be cached. If the Endpoint is a regular expression, IsRegExp must be set.
// CacheTable will then compile it when needed. CacheTable will panic if the regular expression is invalid.
type CacheTableEntry struct {
	Endpoint       string
	IsRegExp       bool
	Expiry         time.Duration
	compiledRegExp *regexp.Regexp
}

// var CacheEverything []CacheTableEntry

func (c *CacheTable) shouldCache(endpoint string) (match bool, expiry time.Duration) {
	c.compileIfNeeded()

	if len(c.Table) == 0 {
		return true, 0
	}

	for _, entry := range c.Table {
		if entry.IsRegExp {
			match = entry.compiledRegExp.MatchString(endpoint)
		} else {
			match = entry.Endpoint == endpoint
		}
		if match {
			expiry = entry.Expiry
			break
		}
	}
	return
}

func (c *CacheTable) compileIfNeeded() {
	if c.compiled {
		return
	}
	var err error
	for index := range c.Table {
		if c.Table[index].IsRegExp {
			c.Table[index].compiledRegExp, err = regexp.Compile(c.Table[index].Endpoint)
			if err != nil {
				panic(fmt.Errorf("cacheTable: invalid regexp '%s': %w", c.Table[index].Endpoint, err))
			}
		}
	}
	c.compiled = true
}
