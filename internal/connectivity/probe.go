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
	Client
}

// NewProbe creates a new Probe
func NewProbe(proxy *url.URL, token string) *Probe {
	return &Probe{
		Client{
			Client: &http.Client{
				Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}},
			Token: token,
		},
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() {
	err := probe.getResponse()
	connected := float64(0)
	if err == nil {
		connected = float64(1)
	}

	metrics.OpenVPNClientStatus.Set(connected)
}

func (probe *Probe) getResponse() error {
	var (
		err      error
		resp     []byte
		response struct {
			IP       string
			Hostname string
			City     string
			Region   string
			Country  string
			Loc      string
			Org      string
			Postal   string
			Timezone string
		}
	)

	if resp, err = probe.call(); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&response)
	}

	log.WithFields(log.Fields{"err": err, "response": response}).Debug("connectivity getResponse")

	return err
}
