package iplocator

import "time"

const cacheDuration = 24 * time.Hour

type cache struct {
	contents map[string]cacheEntry
}

type cacheEntry struct {
	response ipAPIResponse
	expires  time.Time
}

func (c *cache) Add(address string, response ipAPIResponse) {
	if c.contents == nil {
		c.contents = make(map[string]cacheEntry)
	}
	c.contents[address] = cacheEntry{
		response: response,
		expires:  time.Now().Add(cacheDuration),
	}
}

func (c cache) Get(address string) (response ipAPIResponse, found bool) {
	if c.contents == nil {
		return
	}
	var entry cacheEntry
	entry, found = c.contents[address]

	if found {
		if time.Now().After(entry.expires) {
			found = false
		} else {
			response = entry.response
		}
	}

	return
}
