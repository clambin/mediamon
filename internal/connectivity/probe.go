package connectivity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/clambin/mediamon/internal/metrics"
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
				Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
				Timeout:   time.Second * 10,
			},
			Token: token,
		},
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() (err error) {
	var connected float64

	if err = probe.getResponse(); err == nil {
		connected = 1
	} else {
		log.WithField("err", err).Warning("connectivity probe failed")
	}

	metrics.OpenVPNClientStatus.Set(connected)

	return
}

func (probe *Probe) getResponse() (err error) {
	var (
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

	return
}
