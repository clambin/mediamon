package iplocator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client finds the geographic coordinates of an IP address.  It uses https://ip-api.com to look an IP address' location.
type Client struct {
	httpClient *http.Client
	url        string
}

// New creates a new Client
func New(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		url:        ipAPIURL,
	}
}

const ipAPIURL = "http://ip-api.com"

// Locate returns the Location of the specified IP address. No internal validation of the provided IP address is done.
// This is left up entirely to the underlying API.
func (c Client) Locate(address string) (Location, error) {
	var response Location
	resp, err := c.httpClient.Get(c.url + "/json/" + address)
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
	if err != nil {
		return response, fmt.Errorf("invalid response: %w", err)
	}
	if response.Status != "success" {
		err = fmt.Errorf("ip locate failed: %s", response.Status)
	}
	return response, err
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
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}
