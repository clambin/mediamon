package iplocator

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"codeberg.org/clambin/go-common/cache"
)

const ipAPIURL = "http://ip-api.com"

// Client finds the geographic coordinates of an IP address.  It uses https://ip-api.com to look up an IP address' location.
type Client struct {
	cache      *cache.Cache[string, Location]
	httpClient *http.Client
	url        string
}

// New creates a new Client
func New(httpClient *http.Client) *Client {
	return &Client{
		cache:      cache.New[string, Location](time.Hour, 5*time.Minute),
		httpClient: cmp.Or(httpClient, http.DefaultClient),
		url:        ipAPIURL,
	}
}

// Locate returns the Location of the specified IP address. No validation of the provided IP address is done.
// It's left up entirely to the calling client.
func (c Client) Locate(address string) (location Location, err error) {
	// try the cache first
	var ok bool
	if location, ok = c.cache.Get(address); !ok {
		// no cached item. call the remote service.
		if location, err = c.locate(address); err == nil {
			// success! cache the result
			c.cache.Add(address, location)
		}
	}
	return location, err
}

func (c Client) locate(address string) (Location, error) {
	resp, err := c.httpClient.Get(c.url + "/json/" + address)
	if err != nil {
		return Location{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("http: %d - %s", resp.StatusCode, resp.Status)
	}

	var response Location
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Location{}, fmt.Errorf("json: %w", err)
	}
	if response.Status != "success" {
		return Location{}, fmt.Errorf("ip locate failed: %s", response.Status)
	}
	return response, nil
}

type Location struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}
