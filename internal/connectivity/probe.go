package connectivity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"mediamon/internal/metrics"
)

// Probe to measure Plex metrics
type Probe struct {
	client *Client
}

// NewProbe creates a new Probe
func NewProbe(proxy *url.URL, token string) *Probe {
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	return NewProbeWithHTTPClient(client, token)
}

// NewProbeWithHTTPClient creates a probe with a specified http.Client
// Used to stub API calls during unit testing
func NewProbeWithHTTPClient(client *http.Client, token string) *Probe {
	return &Probe{
		client: NewAPIWithHTTPClient(client, token),
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() {
	err := probe.getResponse()
	connected := float64(1)
	if err != nil {
		log.Warningf("%s", err.Error())
		connected = float64(0)
	}

	metrics.OpenVPNClientStatus.Set(connected)
}

func (probe *Probe) getResponse() error {
	var response = struct {
		IP       string
		Hostname string
		City     string
		Region   string
		Country  string
		Loc      string
		Org      string
		Postal   string
		Timezone string
	}{}

	resp, err := probe.client.Call()
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&response)
		if err == nil {
			log.Debug(response)
			return nil
		}
	}
	return err
}
