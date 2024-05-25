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

// Locate finds the longitude and latitude of the specified IP address. No internal validation of the provided IP address is done.
// This is left up entirely to the underlying API.
func (c Client) Locate(ipAddress string) (float64, float64, error) {
	response, err := c.lookup(ipAddress)
	if err != nil {
		return 0, 0, fmt.Errorf("ip locate failed: %w", err)
	}
	if response.Status != "success" {
		return 0, 0, fmt.Errorf("ip locate failed: %s", response.Message)
	}

	return response.Lon, response.Lat, err
}

func (c Client) lookup(ipAddress string) (ipAPIResponse, error) {
	var response ipAPIResponse
	resp, err := c.httpClient.Get(c.url + "/json/" + ipAddress)
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
