package caller

import (
	"bufio"
	"bytes"
	"github.com/clambin/cache"
	"net/http"
	"net/http/httputil"
	"time"
)

// Cacher implements the Caller interface. It will cache calls based in the provided CacheTable
type Cacher struct {
	Client
	Table CacheTable
	Cache cache.Cacher[string, []byte]
}

var _ Caller = &Cacher{}

// NewCacher creates a new Cacher
func NewCacher(httpClient *http.Client, application string, options Options, cacheEntries []CacheTableEntry, cacheExpiry, cacheCleanup time.Duration) *Cacher {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Cacher{
		Client: Client{
			HTTPClient:  httpClient,
			Application: application,
			Options:     options,
		},
		Table: CacheTable{Table: cacheEntries},
		Cache: cache.New[string, []byte](cacheExpiry, cacheCleanup),
	}
}

// Do sends the request and caches the response for future use.
// If a (non-expired) cached response exists for the request's URL, it is returned instead.
//
// Note: only the request's URL is used to find a cached version. Currently, it does not consider the request's method (i.e. GET/PUT/etc).
func (c *Cacher) Do(req *http.Request) (resp *http.Response, err error) {
	key := cacheKey(req)
	body, found := c.Cache.Get(key)
	if found {
		return cachedResponse(body, req)
	}

	resp, err = c.Client.Do(req)

	if err != nil {
		return
	}

	shouldCache, expiry := c.shouldCache(req)
	if !shouldCache {
		return
	}

	var buf []byte
	buf, err = httputil.DumpResponse(resp, true)

	if err == nil {
		c.Cache.AddWithExpiry(key, buf, expiry)
	}
	return
}

func (c Cacher) shouldCache(r *http.Request) (cache bool, expiry time.Duration) {
	cache, expiry = c.Table.shouldCache(r.URL.Path)
	if cache && expiry == 0 {
		expiry = c.Cache.GetDefaultExpiration()
	}
	return
}

func cacheKey(r *http.Request) string {
	return r.URL.String()
}

func cachedResponse(b []byte, r *http.Request) (resp *http.Response, err error) {
	buf := bytes.NewBuffer(b)
	return http.ReadResponse(bufio.NewReader(buf), r)
}
