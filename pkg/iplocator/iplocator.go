package iplocator

import (
	"encoding/json"
	"fmt"
	"github.com/clambin/go-common/cache"
	"github.com/clambin/go-common/httpclient"
	"net/http"
	"time"
)

// Locator finds the geographic coordinates of an IP address
//
//go:generate mockery --name Locator
type Locator interface {
	Locate(ipAddress string) (lon, lat float64, err error)
}

// Client finds the geographic coordinates of an IP address.  It uses https://ip-api.com to look an IP address' location.
type Client struct {
	HTTPClient *http.Client
	URL        string
	ipCache    cache.Cacher[string, ipAPIResponse]
}

// New creates a new Client
func New() *Client {
	return &Client{
		HTTPClient: &http.Client{Transport: httpclient.NewRoundTripper(httpclient.WithCache{
			DefaultExpiry:   24 * time.Hour,
			CleanupInterval: 36 * time.Hour,
		})},
		URL:     ipAPIURL,
		ipCache: cache.New[string, ipAPIResponse](24*time.Hour, 36*time.Hour),
	}
}

var _ Locator = &Client{}

const ipAPIURL = "http://ip-api.com"

// Locate finds the longitude and latitude of the specified IP address. No internal validation of the provided IP address is done.
// This is left up entirely to the underlying API.
func (c *Client) Locate(ipAddress string) (lon, lat float64, err error) {
	response, found := c.ipCache.Get(ipAddress)
	if found && response.Status == "success" {
		lon, lat = response.Lon, response.Lat
		return
	}

	response, err = c.lookup(ipAddress)
	if err != nil {
		err = fmt.Errorf("ip locate failed: %w", err)
		return
	}
	if response.Status != "success" {
		err = fmt.Errorf("ip locate failed: %s", response.Message)
		return
	}

	c.ipCache.Add(ipAddress, response)
	return response.Lon, response.Lat, err
}

func (c *Client) lookup(ipAddress string) (ipAPIResponse, error) {
	var response ipAPIResponse
	resp, err := http.Get(c.URL + "/json/" + ipAddress)
	if err != nil {
		return response, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("ip locate failed: %s", response.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

type ipAPIResponse struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}
